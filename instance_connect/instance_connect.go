package instance_connect

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
)

type AwsSession struct {
	Session   *session.Session
	Ec2Client *ec2.EC2
	IcClient  *ec2instanceconnect.EC2InstanceConnect
}

type InstanceInfo struct {
	ID        string
	Az        string
	PublicDNS string
	User      string
}

func New(profile string, region string) (*AwsSession, error) {
	sess := AwsSession{}
	var err error
	sess.Session, err = session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		return nil, err
	}
	sess.Ec2Client = ec2.New(sess.Session)
	sess.IcClient = ec2instanceconnect.New(sess.Session)
	return &sess, nil
}

func (sess *AwsSession) getInstanceIDFromName(instanceName string, user string) *InstanceInfo {

	filters := []*ec2.Filter{
		{Name: aws.String("tag:Name"), Values: []*string{aws.String(instanceName)}},
	}

	params := &ec2.DescribeInstancesInput{Filters: filters}

	instancesDescription, err := sess.Ec2Client.DescribeInstances(params)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return nil
	}

	instanceInfo := retrieveInstanceInfo(instancesDescription)

	instance := InstanceInfo{
		ID:        *instanceInfo.InstanceId,
		Az:        *instanceInfo.Placement.AvailabilityZone,
		PublicDNS: *instanceInfo.PublicDnsName,
		User:      user,
	}

	return &instance
}

func retrieveInstanceInfo(instanceDescription *ec2.DescribeInstancesOutput) *ec2.Instance {

	if instanceDescription == nil {
		log.Println("Instance description cannot be empty!")
		return nil
	}

	availableInstances := make([]*ec2.Instance, 0, len(instanceDescription.Reservations))

	for _, reservation := range instanceDescription.Reservations {
		availableInstances = append(availableInstances, reservation.Instances...)
	}

	if len(availableInstances) > 1 {
		log.Println("More than one instance found, using the first instance.")
	}

	return availableInstances[0]
}

func (sess *AwsSession) sendSSHKeyToInstanceConnect(instance *InstanceInfo, sshKey string) (*ec2instanceconnect.SendSSHPublicKeyOutput, error) {

	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(instance.Az),
		InstanceId:       aws.String(instance.ID),
		InstanceOSUser:   aws.String(instance.User),
		SSHPublicKey:     aws.String(sshKey),
	}

	result, err := sess.IcClient.SendSSHPublicKey(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return result, err
	}
	return result, nil
}

func (sess *AwsSession) Send(instanceName string, sshKey string, user string) (*InstanceInfo, error) {
	instance := sess.getInstanceIDFromName(instanceName, user)
	result, err := sess.sendSSHKeyToInstanceConnect(instance, sshKey)

	if err == nil && *result.Success {
		log.Println("Key has been successfully added!")
		return instance, err

	}

	log.Println("Unable to add key")

	return nil, err
}
