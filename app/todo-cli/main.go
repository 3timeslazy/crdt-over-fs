package main

import (
	"fmt"
	"os"

	"github.com/3timeslazy/crdt-over-fs/app/todo-cli/tasks"
	"github.com/3timeslazy/crdt-over-fs/app/todo-cli/ui"

	"github.com/3timeslazy/crdt-over-fs/sync"
	"github.com/3timeslazy/crdt-over-fs/sync/crdt/automerge"
	"github.com/3timeslazy/crdt-over-fs/sync/fs/local"
	"github.com/3timeslazy/crdt-over-fs/sync/fs/s3"

	"github.com/aws/aws-sdk-go/aws"
	awscred "github.com/aws/aws-sdk-go/aws/credentials"
	awssess "github.com/aws/aws-sdk-go/aws/session"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jessevdk/go-flags"
	"github.com/pelletier/go-toml"
)

type AppOptions struct {
	Device     string `short:"d" long:"device" required:"true" description:"Device ID"`
	User       string `short:"u" long:"user" required:"true" description:"User"`
	ConfigFile string `short:"c" long:"config" default:"config.toml"`
}

type AppConfig struct {
	FSType string `toml:"fs_type"`

	Local *struct {
		RootDir string `toml:"root_dir"`
	} `toml:"local,omitempty"`

	S3 *struct {
		KeyID     string `toml:"key_id"`
		KeySecret string `toml:"key_secret"`
		Bucket    string `toml:"bucket"`
		Endpoint  string `toml:"endpoint"`
		Region    string `toml:"region"`
	} `toml:"s3,omitempty"`
}

func main() {
	opts := AppOptions{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	confFile, err := os.ReadFile(opts.ConfigFile)
	if err != nil {
		panic(err)
	}
	conf := AppConfig{}
	err = toml.Unmarshal(confFile, &conf)
	if err != nil {
		panic(err)
	}

	var syncCtrl *sync.Controller
	stateID := fmt.Sprintf("%s.%s", opts.Device, opts.User)

	switch conf.FSType {
	case "local":
		syncCtrl = sync.NewController(
			local.NewFS(),
			&automerge.Automerge{},
			stateID,
			conf.Local.RootDir,
		)

	case "s3":
		s3client := newS3Client(conf)
		syncCtrl = sync.NewController(
			s3.NewFS(s3client, conf.S3.Bucket),
			&automerge.Automerge{},
			stateID,
			".",
		)

	default:
		panic(fmt.Sprintf("unknown fs_type %q", conf.FSType))
	}

	err = syncCtrl.InitRootDir()
	if err != nil {
		panic(err)
	}

	app := ui.NewApp(
		opts.Device,
		opts.User,
		tasks.NewManager(syncCtrl),
	)
	prog := tea.NewProgram(app)
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func newS3Client(conf AppConfig) *awss3.S3 {
	creds := awscred.NewStaticCredentials(
		conf.S3.KeyID,
		conf.S3.KeySecret,
		"",
	)
	s3conf := &aws.Config{
		Credentials:      creds,
		Endpoint:         aws.String(conf.S3.Endpoint),
		Region:           aws.String(conf.S3.Region),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := awssess.NewSession(s3conf)
	if err != nil {
		panic(err)
	}

	return awss3.New(sess)
}
