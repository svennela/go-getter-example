package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"errors"
	"strings"
	"time"
	"cloud.google.com/go/storage"
	"net/url"
	"github.com/hashicorp/go-getter"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// gsGetter is a go-getter that works on Cloud Storage using the broker's
// service account. It's incomplete in that it doesn't support directories.
type gsGetter struct{}
// ClientMode is unsupported for gsGetter.
func (g *gsGetter) ClientMode(u *url.URL) (getter.ClientMode, error) {
	return getter.ClientModeInvalid, errors.New("mode is not supported for this client")
}

// Get clones a remote destination to a local directory.
func (g *gsGetter) Get(dst string, u *url.URL) error {
	return errors.New("getting directories is not supported for this client")
}

// GetFile downloads the give URL into the given path. The URL must
// reference a single file. If possible, the Getter should check if
// the remote end contains the same file and no-op this operation.
func (g *gsGetter) GetFile(dst string, u *url.URL) error {
	fmt.Println("---GetFileGetFile---")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client, err := g.client(ctx)
	if err != nil {
		return err
	}

	reader, err := g.objectAt(client, u).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("couldn't open object at %q: %v", u.String(), err)
	}

	return Copy(FromReadCloser(reader), ToFile(dst))
}

func (g *gsGetter) SetClient(_ *getter.Client) {

}

func (gsGetter) objectAt(client *storage.Client, u *url.URL) *storage.ObjectHandle {
	return client.Bucket(u.Hostname()).Object(strings.TrimPrefix(u.Path, "/"))
}

func (gsGetter) client(ctx context.Context) (*storage.Client, error) {
	fmt.Println("---clientclientclientclient---")
	creds, err := google.CredentialsFromJSON(ctx, []byte(GetServiceAccountJson()), storage.ScopeReadOnly)
	if err != nil {
		return nil, errors.New("couldn't get JSON credentials from the enviornment")
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(creds), option.WithUserAgent(CustomUserAgent))
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to Cloud Storage: %v", err)
	}
	return client, nil
}


var gitgetters = map[string]getter.Getter{
	"git":   new(getter.GitGetter),
	"http":  new(getter.GitGetter), // this supports git downloads through https/http
	"https": new(getter.GitGetter), // this supports git downloads through https/http
	//	"http":  getterHTTPGetter, to support actual http downloads
	//  "https": getterHTTPGetter, to support actual https downloads
	"gcs":    new(getter.GCSGetter),
	"gs":    &gsGetter{},
	//"gs":    new(getter.GCSGetter),
	"s3":    new(getter.S3Getter),
}

func main(){
	fmt.Println("--------")
	dest := os.Args[1]
	src := os.Args[2]
	download(src,dest)
	fmt.Println("--------")
}	

func download(src,dest string){
	log.Printf("fetching %q to %q", src, dest)
	client := &getter.Client{
		Ctx:     context.Background(),
		Src:     src,
		Dst:     dest,
		Mode:    getter.ClientModeDir,
		Options: []getter.ClientOption{},
		Getters: gitgetters,
	}
	//download the files
	if err := client.Get(); err != nil {
		fmt.Fprintf(os.Stderr, "Error getting path %s: %v", client.Src, err)
		os.Exit(1)
	}
}

