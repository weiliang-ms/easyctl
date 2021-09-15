package export

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"os"
)

func LocalImageList() {

	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("连接docker失败 -> %s", err.Error())
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {
		panic(err)
	}

	imageListPath := "local-image-list.txt"

	if _, err := os.Stat(imageListPath); err == nil {
		_ = os.Remove(imageListPath)
	}

	imageList, _ := os.OpenFile(imageListPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)

	log.Println("########## 镜像列表 ##########")
	for _, v := range images {
		for _, t := range v.RepoTags {
			fmt.Println(t)
			_, _ = imageList.WriteString(t)
			_, _ = imageList.WriteString("\n")
		}
	}
	log.Println("########## 镜像列表 ##########")

	defer imageList.Close()
}
