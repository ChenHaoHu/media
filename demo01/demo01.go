package main

import (
	"fmt"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/cgo/ffmpeg"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/flv"
	"image/png"
	"os"
)

func main() {
	// 注册视频格式
	format.RegisterAll()

	fmt.Println("----------------------")
	// 打开视频文件
	file, err := os.Open("source.200kbps.768x320.flv")
	if err != nil {
		fmt.Println("Error opening video file:", err)
		return
	}
	defer file.Close()

	// 创建视频解码器
	demuxer := flv.NewDemuxer(file)

	streams, err := demuxer.Streams()
	if err != nil {
		fmt.Println("demuxer.Streams err", err)
		return
	}
	fmt.Println("streams count:", len(streams))
	for _, stream := range streams {
		if stream.Type().IsAudio() {
			astream := stream.(av.AudioCodecData)
			fmt.Println("audio stream: ", stream.Type(), astream.SampleRate(), astream.SampleFormat(), astream.ChannelLayout())
		} else if stream.Type().IsVideo() {
			vstream := stream.(av.VideoCodecData)
			fmt.Println("video stream: ", vstream.Type(), vstream.Width(), vstream.Height())
			videoDecoder, err := ffmpeg.NewVideoDecoder(vstream)
			if err != nil {
				fmt.Println("ffmpeg.NewVideoDecoder err:", err)
				return
			}
			// 读取视频帧
			for {
				packet, err := demuxer.ReadPacket()
				if err != nil {
					fmt.Println("demuxer.ReadPacket err:", err)
					break
				}

				if packet.IsKeyFrame {
					fmt.Println(packet.CompositionTime)
					frame, err := videoDecoder.Decode(packet.Data)
					if err != nil {
						fmt.Println("decode err:", err)
						return
					}
					if frame == nil {
						continue
					}
					err = saveFrameAsImage(frame)
					if err != nil {
						fmt.Println("saveFrameAsImage err:", err)
					}
					//break // 只处理第一帧
				}
				fmt.Println("next packet")
			}
		}
	}

}

var index = 0

func saveFrameAsImage(frame *ffmpeg.VideoFrame) error {

	subImage := frame.Image.SubImage(frame.Image.Bounds())

	// 创建图像文件
	file, err := os.Create(fmt.Sprintf("output/%d.png", index))
	if err != nil {
		return err
	}
	defer file.Close()

	index++
	// 将图像保存为PNG文件
	return png.Encode(file, subImage)
}
