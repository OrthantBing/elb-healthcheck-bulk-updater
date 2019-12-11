package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func main() {
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")
	UnhealthyThresholdCount := 2
	HealthCheckIntervalSeconds := 7
	lbResource := "arn:aws:elasticloadbalancing:ap-south-1:379434283010:loadbalancer/app/Staging/42db6c66a102a70f"
	token := ""

	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecret, token)
	cfg := aws.NewConfig().WithRegion(awsRegion).WithCredentials(creds)

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: *cfg,
	})
	if err != nil {
		fmt.Println("Unable to create session", err)
		os.Exit(1)
	}

	// svc := s3.New(sess)
	// result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	// if err != nil {
	// 	log.Println("Failed to list buckets", err)
	// 	return
	// }

	// log.Println("Buckets:")

	// for _, bucket := range result.Buckets {
	// 	log.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
	// }
	svc := elbv2.New(sess)
	input := &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: aws.String(lbResource),
	}

	result, err := svc.DescribeTargetGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elbv2.ErrCodeLoadBalancerNotFoundException:
				fmt.Println(elbv2.ErrCodeLoadBalancerNotFoundException, aerr.Error())
			case elbv2.ErrCodeTargetGroupNotFoundException:
				fmt.Println(elbv2.ErrCodeTargetGroupNotFoundException, aerr.Error())
			case elbv2.ErrCodeListenerNotFoundException:
				fmt.Println(elbv2.ErrCodeListenerNotFoundException, aerr.Error())
			case elbv2.ErrCodeRuleNotFoundException:
				fmt.Println(elbv2.ErrCodeRuleNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	for _, val := range result.TargetGroups {
		fmt.Println("TargetGroupName: ", *val.TargetGroupName)

		fmt.Println("Current HealthCheckIntervalSeconds: ", *val.HealthCheckIntervalSeconds)
		fmt.Println("Current UnhealthyThresholdCount: ", *val.UnhealthyThresholdCount)
		fmt.Println("Current HealthCheckTimeoutSeconds", *val.HealthCheckTimeoutSeconds)

		fmt.Println()
		fmt.Println("Change HealthCheckIntervalSeconds", HealthCheckIntervalSeconds)
		fmt.Println("Change UnhealthyThresholdCount", UnhealthyThresholdCount)

		if *val.HealthCheckIntervalSeconds == int64(HealthCheckIntervalSeconds) &&
			*val.UnhealthyThresholdCount == int64(UnhealthyThresholdCount) {
			continue
		}

		fmt.Print("Confirm Change (x to exit) y/n/x: ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		fmt.Println()
		if strings.ToLower(input.Text()) == "y" {
			input := &elbv2.ModifyTargetGroupInput{
				HealthCheckIntervalSeconds: aws.Int64(int64(HealthCheckIntervalSeconds)),
				UnhealthyThresholdCount:    aws.Int64(int64(UnhealthyThresholdCount)),
				TargetGroupArn:             val.TargetGroupArn,
			}

			result, err := svc.ModifyTargetGroup(input)
			if err != nil {
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case elbv2.ErrCodeTargetGroupNotFoundException:
						fmt.Println(elbv2.ErrCodeTargetGroupNotFoundException, aerr.Error())
					case elbv2.ErrCodeInvalidConfigurationRequestException:
						fmt.Println(elbv2.ErrCodeInvalidConfigurationRequestException, aerr.Error())
					default:
						fmt.Println(aerr.Error())
					}
				} else {
					// Print the error, cast err to awserr.Error to get the Code and
					// Message from an error.
					fmt.Println(err.Error())
				}
				return
			}

			fmt.Println(result)
		} else if strings.ToLower(input.Text()) == "x" {
			fmt.Println("Exiting")
			os.Exit(0)
		} else {
			fmt.Println(val.TargetGroupName, " Unchanged")
			fmt.Println()
		}
	}
}
