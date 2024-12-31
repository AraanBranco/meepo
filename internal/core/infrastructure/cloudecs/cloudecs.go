package cloudecs

import (
	"context"
	"fmt"

	"github.com/AraanBranco/meepow/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"go.uber.org/zap"
)

func LaunchContainer(ecsClient *ecs.Client, config config.Config, referenceID string) (string, error) {
	cluster := config.GetString("providers.aws.clusterName")
	taskDefinition := config.GetString("providers.aws.taskDefinition")
	containerName := config.GetString("providers.aws.containerName")
	logger := zap.L().With(zap.String("service", "ecs"))

	// TODO: think how get the bot username and password (maybe from file?)
	envVars := []types.KeyValuePair{
		{
			Name:  aws.String("MEEPOW_REFERENCE_ID"),
			Value: aws.String(referenceID),
		},
		{
			Name:  aws.String("MEEPOW_BOT_ALLOWCHEATS"),
			Value: aws.String(config.GetString("bot.allowCheats")),
		},
		{
			Name:  aws.String("MEEPOW_BOT_USERNAME"),
			Value: aws.String(config.GetString("bot.username")),
		},
		{
			Name:  aws.String("MEEPOW_BOT_PASSWORD"),
			Value: aws.String(config.GetString("bot.password")),
		},
		{
			Name:  aws.String("MEEPOW_ADAPTERS_REDIS_URI"),
			Value: aws.String(config.GetString("adapters.redis.uri")),
		},
	}

	// Sobrescreve as definições do container na task definition
	overrides := &types.TaskOverride{
		ContainerOverrides: []types.ContainerOverride{
			{
				Name:        aws.String(containerName),
				Environment: envVars,
				Command:     []string{"start", "bot"},
			},
		},
	}

	// Executa a task
	input := &ecs.RunTaskInput{
		Cluster:        aws.String(cluster),
		TaskDefinition: aws.String(taskDefinition),
		Overrides:      overrides,
		LaunchType:     types.LaunchTypeFargate,
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        []string{config.GetString("providers.aws.subnet1"), config.GetString("providers.aws.subnet2")},
				SecurityGroups: []string{config.GetString("providers.aws.securityGroup")},
				AssignPublicIp: types.AssignPublicIpEnabled,
			},
		},
	}

	output, err := ecsClient.RunTask(context.TODO(), input)
	if err != nil {
		logger.Error("Error running task", zap.Error(err))
		return "", err
	}

	logger.Info("Task executed successfully!")
	for _, task := range output.Tasks {
		fmt.Println(*task.TaskArn)
	}
	return *output.Tasks[0].TaskArn, nil
}
