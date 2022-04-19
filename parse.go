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
	tableNumber := 0

	tc := make([]int, 0)
	tn := make([]int, 0)
	l := make([][]int, 16)
	for i := range l {
		l[i] = make([]int, 0)
	}
	vdc := make([][][]int, 16)
	for i := range vdc {
		vdc[i] = make([][]int, 16)
		for j := range vdc[i] {
			vdc[i][j] = make([]int, 0)
		}
	}
	vac := make([][][]int, 16)
	for i := range vac {
		vac[i] = make([][]int, 16)
		for j := range vac[i] {
			vac[i][j] = make([]int, 0)
		}
	}

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

		// DHTセグメント
		case 0xffc4:
			dhtNR := bytes.NewReader(data)
			fmt.Println("Segment : DHT", "Table:", tableNumber)

			lp := true
			for lp {
				len := 19
				first, second := Read4BitsAsInt(dhtNR)
				tc = append(tc, first)
				tn = append(tn, second)
				for i := 0; i < 16; i++ {
					lnum := ReadBytesAsInt(dhtNR, 1)
					l[tableNumber] = append(l[tableNumber], lnum)
					len += lnum
				}
				if tc[tableNumber] == 0 {
				// DC成分
					for i := 0; i < 16; i++ {
						for k := 0; k < l[tableNumber][i]; k++ {
							vdc[tableNumber][i] = append(vdc[tableNumber][i], ReadBytesAsInt(dhtNR, 1))
						}
					}
				} else if tc[tableNumber] == 1 {
					// AC成分
					for i := 0; i < 16; i++ {
						for k := 0; k < l[tableNumber][i]; k++ {
							runlength, databit := Read4BitsAsInt(dhtNR)
							vac[tableNumber][i] = append(vac[tableNumber][i], runlength)
							vac[tableNumber][i] = append(vac[tableNumber][i], databit)
						}
					}
				}

				fmt.Println("Table Class:", tc[tableNumber], "Table Number:", tn[tableNumber])
				if tc[tableNumber] == 0 {
				// DC成分
					for i := 0; i < 16; i++ {
						fmt.Println("Table:", i,"BitNum:",l[tableNumber][i], "DataBitNum:",vdc[tableNumber][i])
					}
				} else if tc[tableNumber] == 1{
				// AC成分
					for i:= 0; i<16; i++{
						fmt.Println("Table", i,"BitNum:",l[tableNumber][i], "RunLengthNum:",vac[tableNumber][i],"DataBitNum:",vac[tableNumber][i])
					}
				}
				tableNumber++

				if len == length {
					lp = false
				}
			}

		// SOF0セグメント
		case 0xffc0:
			sofNR := bytes.NewReader(data)
			fmt.Println("Segment : SOF")
			elementPrecision := ReadBytesAsInt(sofNR, 1)
			heightSize := ReadBytesAsInt(sofNR, 2)
			widthSize := ReadBytesAsInt(sofNR, 2)
			elementNum := ReadBytesAsInt(sofNR, 1)
			fmt.Println("Element Precision:", elementPrecision, "Image Size Height:", heightSize, "Width:", widthSize)
			cn := make([]int, 0)
			hn := make([]int, 0)
			vn := make([]int, 0)
			tqn := make([]int, 0)
			for i := 0; i < elementNum; i++ {
				cn = append(cn, ReadBytesAsInt(sofNR, 1))
				first, second := Read4BitsAsInt(sofNR)
				hn = append(hn, first)
				vn = append(vn, second)
				tqn = append(tqn, ReadBytesAsInt(sofNR, 1))
				fmt.Println("Element No.", i, "ID:", cn[i], "Horizontal Level;", hn[i], "Vertical Level:", vn[i], "Table Number:", tqn[i])
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
