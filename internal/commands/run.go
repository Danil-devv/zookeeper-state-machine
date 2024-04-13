package commands

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"hw/internal/commands/cmdargs"
	"hw/internal/config"
	"hw/internal/depgraph"
	"log/slog"
	"strings"
)

func InitRunCommand(ctx context.Context) (cobra.Command, error) {
	cmdArgs := cmdargs.RunArgs{}
	cmd := cobra.Command{
		Use:   "run",
		Short: "Starts a leader election node",
		Long: `This command starts the leader election node that connects to zookeeper
		and starts to try to acquire leadership by creation of ephemeral node`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dg := depgraph.New()
			logger, err := dg.GetLogger()
			if err != nil {
				return fmt.Errorf("get logger: %w", err)
			}
			logger.Info(
				"args received",
				slog.String("servers", strings.Join(cmdArgs.ZookeeperServers, ", ")),
				slog.Duration("leader-timeout", cmdArgs.LeaderTimeout),
				slog.Duration("attempter-timeout", cmdArgs.AttempterTimeout),
				slog.Duration("failover-timeout", cmdArgs.AttempterTimeout),
				slog.String("filedir", cmdArgs.FileDir),
				slog.Int("storage-capacity", cmdArgs.StorageCapacity),
				slog.Int("failover-attempts-count", cmdArgs.FailoverAttemptsCount),
			)

			runner, err := dg.GetRunner()
			if err != nil {
				return fmt.Errorf("get runner: %w", err)
			}
			logger.Debug("successfully get runner")

			conn, err := dg.GetZkConn(&cmdArgs)
			if err != nil {
				return fmt.Errorf("get zk connection: %w", err)
			}
			logger.Debug("successfully get zk connection")

			b, err := dg.GetBasicState(&cmdArgs, logger, conn)
			if err != nil {
				return fmt.Errorf("get basic state: %w", err)
			}
			logger.Debug("successfully get basic state")

			initState, err := dg.GetInitState(b)
			if err != nil {
				return fmt.Errorf("get init state: %w", err)
			}
			logger.Debug("successfully get init state")

			logger.Info("starting state machine")
			err = runner.Run(ctx, initState, b)
			if err != nil {
				return fmt.Errorf("run states: %w", err)
			}
			return nil
		},
	}

	setCmdArgs(&cmd, &cmdArgs)

	return cmd, nil
}

func setCmdArgs(cmd *cobra.Command, cmdArgs *cmdargs.RunArgs) {
	conf := config.GetEnvConfig()
	cmd.Flags().StringSliceVarP(
		&(cmdArgs.ZookeeperServers),
		"zk-servers",
		"s",
		conf.ZookeeperServers,
		"Set the zookeeper servers.",
	)
	cmdArgs.LeaderTimeout = *cmd.Flags().Duration(
		"leader-timeout",
		conf.LeaderTimeout,
		"Sets the frequency at which the leader writes the file to disk.",
	)
	cmdArgs.AttempterTimeout = *cmd.Flags().Duration(
		"attempter-timeout",
		conf.AttempterTimeout,
		"Sets the frequency with which the attempter tries to become a leader.",
	)
	cmdArgs.FailoverTimeout = *cmd.Flags().Duration(
		"failover-timeout",
		conf.FailoverTimeout,
		"Sets the frequency with which the failover tries to resume its work.",
	)
	cmdArgs.FileDir = *cmd.Flags().String(
		"file-dir",
		conf.FileDir,
		"Sets the directory where the leader must write files.",
	)
	cmdArgs.StorageCapacity = *cmd.Flags().Int(
		"storage-capacity",
		conf.StorageCapacity,
		"Sets the maximum number of files in the file-dir directory.",
	)
	cmdArgs.FailoverAttemptsCount = *cmd.Flags().Int(
		"attempts-count",
		conf.FailoverAttemptsCount,
		"Sets the maximum number of failover attempts",
	)
}
