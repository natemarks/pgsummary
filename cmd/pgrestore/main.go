package main

import (
	"flag"
	"github.com/natemarks/easyaws/rds"
	"github.com/natemarks/pgsummary/util"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	delete_after := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	logger := log.With().Str("test_key", "test_value").Logger()
	instancePtr := flag.String("instance", "", "Postgres RDS instance name")
	flag.Parse()
	restoreInstance := "deleteme-" + *instancePtr
	region, err := util.GetAWSRegionEnvVar()
	util.CheckError(err, &logger)

	logger.Info().Msgf("Using AWS_DEFAULT_REGION: %s", region)

	snapshotId, err := rds.GetLatestSnapshotId(*instancePtr, &logger)
	util.CheckError(err, &logger)

	dbSubnetgroupName, err := rds.GetSubnetGroup(*instancePtr, &logger)
	util.CheckError(err, &logger)

	vpcSecurityGroupIDs, err := rds.GetVPCSecurityGroups(*instancePtr, &logger)
	tag := []rds.Tag{
		{Key: "deleteme", Value: "true"},
		{Key: "deleteme_after", Value: delete_after},
	}
	rsInput := rds.RestoreSnapshotIdInput{
		DBInstanceIdentifier: restoreInstance,
		DBSnapshotIdentifier: snapshotId,
		DBSubnetGroupName:    dbSubnetgroupName,
		VpcSecurityGroupIds:  vpcSecurityGroupIDs,
		Tags:                 tag,
	}
	_, err = rds.RestorePGSnapshotId(rsInput, &logger)
	util.CheckError(err, &logger)

	log.Info().Msgf("Restoring snapshot: %s to instance: %s", snapshotId, restoreInstance)
}
