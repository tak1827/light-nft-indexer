package nftcli

import (

	// "github.com/tak1827/light-nft-indexer/apiclient"
	// "github.com/tak1827/light-nft-indexer/data"
	// job "github.com/tak1827/light-nft-indexer/job/fetch"
	"github.com/tak1827/light-nft-indexer/log"
	// "github.com/tak1827/light-nft-indexer/store"
	// "github.com/aws/aws-sdk-go/aws"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "indexing by fetching logs",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLogColor(logWithColor)

		getConfig()
		getAwsConfig()
		getS3Config()

		runfetch()
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

func runfetch() {
	// var (
	// 	ctx    = context.Background()
	// 	db     = store.NewDynamoDB(store.TableName, store.WithProfileOpt(AwsProfile), store.WithRegionOpt(AwsRegion), store.WithEndpointOpt(DynamoEndpoint))
	// 	batch  = db.Batch()
	// 	s3conf = aws.NewConfig().WithRegion(AwsRegion).WithLogLevel(aws.LogOff)
	// 	now    = time.Now()

	// 	fetchs []*data.fetch
	// )

	// client, err := apiclient.NewS3Client(AwsProfile, S3BucketName, s3conf)
	// if err != nil {
	// 	logger.Fatal().Msgf("failed to initialize s3 client: %d", err)
	// }

	// // 実行すべき配当をリストアップ
	// if err := db.ListEntityWithStatusAndFilterByTimestampLE(ctx, data.Typefetch.String(), data.fetchScheduled.String(), "payment_date", &now, &fetchs); err != nil || fetchs == nil {
	// 	logger.Info().Msgf("%v", fetchs)
	// 	if err != nil && err != store.ErrNotFound {
	// 		logger.Fatal().Msgf("failed to list fetchs: %d", err)
	// 	}
	// 	logger.Info().Msgf("no fetchs detected. Job ends at: %v", now)
	// 	return
	// }

	// var items []interface{}
	// for _, fetch := range fetchs {
	// 	logger.Info().Msgf("fetch process started with: %v", fetch)

	// 	snapshotColumns, err := job.ReadSnapshotColumns(&client, fetch.GetSnapshotObjectKey())
	// 	if err != nil {
	// 		logger.Fatal().Msgf("failed to read snapshot csv: %s", fetch.GetSnapshotObjectKey())
	// 	}

	// 	for _, col := range snapshotColumns {
	// 		userfetch := data.NewUserfetch(col.fetchId, col.InvestorId, fetch.GetTokenId(), uint32(col.fetchAmount), &now)
	// 		items = append(items, userfetch)
	// 	}

	// 	fetch.Status = data.fetchCompleted.String()
	// 	fetch.StatusWithCreatedAt = data.BuildCompositeKey(fetch.GetStatus(), fetch.GetCreatedAt().String())
	// 	fetch.UpdatedAt = &now

	// 	items = append(items, fetch)

	// 	logger.Info().Msgf("fetch process succeeded: %v", fetch)
	// }

	// for i := 0; i < len(items); i++ {
	// 	switch item := items[i].(type) {
	// 	case data.Userfetch:
	// 		batch.Put(item)
	// 	case *data.fetch:
	// 		batch.Update(*item)
	// 	}
	// 	if i%25 == 0 || i == len(items)-1 {
	// 		if err = batch.Commit(ctx); err != nil {
	// 			logger.Fatal().Msgf("failed to commit batch: %d", err)
	// 		}
	// 	}
	// }
}
