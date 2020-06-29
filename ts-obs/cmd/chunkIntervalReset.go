package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// chunkIntervalResetCmd represents the chunk-interval reset command
var chunkIntervalResetCmd = &cobra.Command{
	Use:   "reset <metric>",
	Short: "Resets the chunk interval for a specific metric back to the default",
	Args:  cobra.ExactArgs(1),
	RunE:  chunkIntervalReset,
}

func init() {
	chunkIntervalCmd.AddCommand(chunkIntervalResetCmd)
	chunkIntervalResetCmd.Flags().StringP("user", "U", "postgres", "database user name")
	chunkIntervalResetCmd.Flags().StringP("dbname", "d", "postgres", "database name to connect to")
}

func chunkIntervalReset(cmd *cobra.Command, args []string) error {
	var err error

	metric := args[0]

	var user string
	user, err = cmd.Flags().GetString("user")
	if err != nil {
		return fmt.Errorf("could not reset chunk interval for %v: %w", metric, err)
	}

	var dbname string
	dbname, err = cmd.Flags().GetString("dbname")
	if err != nil {
		return fmt.Errorf("could not reset chunk interval for %v: %w", metric, err)
	}

	pool, err := OpenConnectionToDB(namespace, name, user, dbname, FORWARD_PORT_TSDB)
	if err != nil {
		return fmt.Errorf("could not reset chunk interval for %v: %w", metric, err)
	}
	defer pool.Close()

	fmt.Printf("Resetting chunk interval for %v back to default\n", metric)
	_, err = pool.Exec(context.Background(), "SELECT prom_api.reset_metric_chunk_interval($1)", metric)
	if err != nil {
		return fmt.Errorf("could not reset chunk interval for %v: %w", metric, err)
	}

	return nil
}