package aws


import (
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


type FirehoseGenerator struct {
	AWSService
}

// createResources iterate on all buckets
// for each bucket we check region and choose only bucket from set region
// for each bucket try get bucket policy, if policy exist create additional NewTerraformResource for policy
func (g FirehoseGenerator) createResources(sess *session.Session, streamNames []*string, region string) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, streamName := range streamNames {
		resourceName := aws.StringValue(streamName)
		// check if bucket in region
		//if s3.NormalizeBucketLocation(aws.StringValue(location.LocationConstraint)) == region {
		resources = append(resources, terraform_utils.NewResource(
			resourceName,
			resourceName,
			"aws_kinesis_firehose_delivery_stream",
			"aws",
			map[string]string{},
			S3AllowEmptyValues,
			map[string]string{}))
	}
	return resources
}

// Generate TerraformResources from AWS API,
// from each s3 bucket create 2 TerraformResource(bucket and bucket policy)
// Need bucket name as ID for terraform resource
func (g *FirehoseGenerator) InitResources() error {
	sess := g.generateSession()
	svc := firehose.New(sess)
	var streamNames []*string
	for {
		output, err := svc.ListDeliveryStreams(&firehose.ListDeliveryStreamsInput{})
		if err != nil {
			return err
		}
		streamNames = append(streamNames, output.DeliveryStreamNames...)
		if *output.HasMoreDeliveryStreams == false {
			break
		}
	}

	g.Resources = g.createResources(sess, streamNames, g.GetArgs()["region"])
	g.PopulateIgnoreKeys()
	return nil
}


//// PostGenerateHook for add bucket policy json as heredoc
//// support only bucket with policy
//func (g *FirehoseGenerator) PostConvertHook() error {
//	for i, resource := range g.Resources {
//		if resource.InstanceInfo.Type != "aws_s3_bucket_policy" {
//			continue
//		}
//		policy := resource.Item["policy"].(string)
//		g.Resources[i].Item["policy"] = fmt.Sprintf(`<<POLICY
//%s
//POLICY`, policy)
//	}
//	return nil
//}

//
//func (g *FirehoseGenerator) ParseFilter(rawFilter []string) {
//	g.Filter = map[string][]string{}
//	for _, resource := range rawFilter {
//		t := strings.Split(resource, "=")
//		if len(t) != 2 {
//			log.Println("Pattern for filter must be resource_type=id1:id2:id4")
//			continue
//		}
//		resourceName, resourcesID := t[0], t[1]
//		g.Filter[resourceName] = strings.Split(resourcesID, ":")
//		if resourceName == "aws_s3_bucket" {
//			g.Filter["aws_s3_bucket_policy"] = strings.Split(resourcesID, ":")
//		}
//	}
//}