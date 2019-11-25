package cmd

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/fox-one/pkg/mq"
	"github.com/fox-one/pkg/mq/sqs"
	"github.com/fox-one/pkg/text/localizer"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/plugin/parser"
	"github.com/yiplee/blockquiz/plugin/shuffler"
	"github.com/yiplee/blockquiz/service/wallet"
	"golang.org/x/text/language"
)

func provideLocalizer() *localizer.Localizer {
	var files []string
	_ = filepath.Walk(cfg.I18n.Path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			files = append(files, path)
		}

		return nil
	})

	return localizer.New(language.Chinese, files...)
}

func provideWalletService() core.WalletService {
	return wallet.New(wallet.Config{
		ClientID:   cfg.Bot.ClientID,
		SessionID:  cfg.Bot.SessionID,
		PinToken:   cfg.Bot.PinToken,
		Pin:        cfg.Bot.Pin,
		SessionKey: cfg.Bot.SessionKey,
	})
}

func provideParser() core.CommandParser {
	return parser.New()
}

func provideShuffler() core.CourseShuffler {
	return shuffler.Rand()
}

func provideAwsSession() *session.Session {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWS.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AWS.Key, cfg.AWS.Secret, ""),
	})

	if err != nil {
		panic(err)
	}

	return s
}

func providePubSub(s *session.Session) mq.PubSub {
	return sqs.New(s, cfg.AWS.QueueURL)
}
