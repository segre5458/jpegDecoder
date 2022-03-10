package main

import (
	"bytes"
	"fmt"
)

func Parse(buffer []byte) (err error) {
	r := bytes.NewReader(buffer)
	CheckSOI(ReadBytesAsInt(r, 2))

	loop := true
	bytelength := 2

	for loop {
		marker := ReadBytesAsInt(r, 2)
		length := ReadBytesAsInt(r, 2)
		data := ReadBytes(r, length-2)
		bytelength += length + 2

		switch marker {
		// APP0セグメント
		case 0xffe0:
			app0NR := bytes.NewReader(data)
			fmt.Println("Segment : APP0")
			jfifsign := ReadBytes(app0NR, 5)
			version := ReadBytesAsInt(app0NR, 2)
			pixelunit := ReadBytesAsInt(app0NR, 1)
			pixelwidth := ReadBytesAsInt(app0NR, 2)
			pixelheight := ReadBytesAsInt(app0NR, 2)
			thumbnailwidth := ReadBytesAsInt(app0NR, 1)
			thumbnailheight := ReadBytesAsInt(app0NR, 1)
			jfif := []byte("JFIF")
			jfif = append(jfif, 0x00)
			if bytes.Equal(jfifsign, jfif) {
				fmt.Println("This file is JPEG")
			} else {
				fmt.Println("This file is not JPEG")
				return
			}
			fmt.Println(version)
			var pixelunitstr string
			switch pixelunit {
			case 1:
				pixelunitstr = "inch"
			case 2:
				pixelunitstr = "cm"
			default:
				pixelunitstr = "undefined"
			}
			fmt.Println("Version:", version, "Pixel Unit:", pixelunitstr, "Pixel Width:", pixelwidth, "Pixel Height:", pixelheight)
			if thumbnailwidth != 0 && thumbnailheight != 0 {
				thumbnaildata := ReadBytes(app0NR, length-16)
				fmt.Println("This file has thumbnail")
				fmt.Println("Thumbnail Width", thumbnailwidth, "Height:", thumbnailheight, "Length:", len(thumbnaildata))
			} else {
				fmt.Println("This file has not thumbnail")
			}

		// DQTセグメント
		case 0xffdb:
			dqtNR := bytes.NewReader(data)
			fmt.Println("Segment : DQT")
			tablePrecision := make([]int, 0)
			tableNum := make([]int, 0)
			table := make([][]byte, 0)
			// ベースライン方式にのみ対応
			for i := 0; i < (length-2)/65; i++ {
				first, second := Read4BitsAsInt(dqtNR)
				tablePrecision = append(tablePrecision, first)
				tableNum = append(tableNum, second)
				table = append(table, ReadBytes(dqtNR, 64))
				fmt.Println("TableNum:", tableNum[i], "Precision:", tablePrecision[i], "table:", table[i])
			}

		// SOF0セグメント
		case 0xffc0:
			sofNR := bytes.NewReader(data)
			fmt.Println("Segment : SOF")
			elementPrecision := ReadBytesAsInt(sofNR, 1)
			heightSize := ReadBytesAsInt(sofNR, 2)
			widthSize := ReadBytesAsInt(sofNR, 2)
			elementNum := ReadBytesAsInt(sofNR, 1)
			fmt.Println("Element Precision:",elementPrecision,"Image Size Height:",heightSize,"Width:",widthSize)
			cn := make([]int,0)
			hn := make([]int,0)
			vn := make([]int,0)
			tqn := make([]int,0)
			for i := 0; i < elementNum; i++ {
				cn = append(cn, ReadBytesAsInt(sofNR,1))
				first,second := Read4BitsAsInt(sofNR)
				hn = append(hn, first)
				vn = append(vn, second)
				tqn = append(tqn, ReadBytesAsInt(sofNR,1))
				fmt.Println("Element No.",i,"ID:",cn[i],"Horizontal Level;",hn[i],"Vertical Level:",vn[i],"Table Number:",tqn[i])
			}

		// SOSセグメント
		case 0xffda:
			sosNR := bytes.NewReader(data)
			fmt.Println("Segment : SOS")
			componentnum := ReadBytesAsInt(sosNR, 1)
			cs := make([]int, 0)
			td := make([]int, 0)
			ta := make([]int, 0)
			for i := 0; i < componentnum; i++ {
				cs = append(cs, ReadBytesAsInt(sosNR, 1))
				first, second := Read4BitsAsInt(sosNR)
				td = append(td, first)
				ta = append(ta, second)
			}
			ss := ReadBytesAsInt(sosNR, 1)
			sr := ReadBytesAsInt(sosNR, 1)
			first, second := Read4BitsAsInt(sosNR)
			ah := first
			al := second
			fmt.Println("Component Number:", componentnum)
			for i := 0; i < componentnum; i++ {
				fmt.Println("Component No", i, "ID:", cs[i], "DC Number:", td[i], "AC Number:", ta[i])
			}
			fmt.Println("Start Number:", ss, "End Number:", sr, "Ah:", ah, "Al:", al)

			loop = false
		}
	}

	imgData := ReadBytes(r, len(buffer)-2-bytelength)
	fmt.Println("imgData:", imgData)
	CheckEOI(ReadBytesAsInt(r, 2))

	return err
}
