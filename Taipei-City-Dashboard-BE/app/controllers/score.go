package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
)

type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   struct {
		Type        string      `json:"type"`
		Coordinates []float64   `json:"coordinates"`
	} `json:"geometry"`
}

func Rad2Deg(radians float64) float64 {
    return radians * 180 / math.Pi
}

func Deg2Rad(degrees float64) float64 {
    return degrees * math.Pi / 180
}

func haversine(longitude1, latitude1, longitude2, latitude2 float64) float64 {
    theta := longitude1 - longitude2

    distance := 60 * 1.1515 * Rad2Deg(
        math.Acos(
            (math.Sin(Deg2Rad(latitude1)) * math.Sin(Deg2Rad(latitude2))) +
            (math.Cos(Deg2Rad(latitude1)) * math.Cos(Deg2Rad(latitude2)) * math.Cos(Deg2Rad(theta))),
        ),
    )

	return math.Round(distance*1.609344*100) / 100
}

// Point 代表一個經緯度的點
type Point struct {
	Longitude, Latitude float64
}

// DistanceToPolygon 計算一個點到一個多邊形的最短距離
func DistanceToPolygon(point Point, polygon []Point) float64 {
	minDistance := math.Inf(1)
	for i := 0; i < len(polygon); i++ {
		target := polygon[i]
		distance := haversine(point.Latitude, point.Longitude, target.Latitude, target.Longitude)
		minDistance = math.Min(minDistance, distance)
	}
	return minDistance
}

/**
	x: value for longitude
	y: value for latitude 
	filename: file name of data
	weight: weight rate e.g. 0.07
	calcRange: calculate range
*/
func GetWeightScore(x float64, y float64, filename string,
	 weight float64, calcRange float64) float64 {
	// 設定資料夾路徑
	datasourceDir := "./datasource"
	// 組合完整路徑
	filePath := filepath.Join(datasourceDir, filename)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("無法讀取檔案: %v", err)
		return 0
	}
	// 解析 JSON 內容到自定義結構體
	var data struct {
		Features []GeoJSONFeature `json:"features"`
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		fmt.Printf("無法解析 JSON: %v", err)
		return 0
	}
	// 計算輸入座標與每個點的距離
    distances := make(map[float64][]float64) // 儲存距離和對應的座標
    for _, feature := range data.Features {
        coords := feature.Geometry.Coordinates
		distance := haversine(x, y, coords[0], coords[1])
        distances[distance] = coords
    }
    // 轉換距離和座標的 map 為 JSON 可接受的格式
    var result = 0.0;
	a := 0.5  // 指數增長速率
    for distance := range distances {
		if (distance < calcRange) {
			distanceScore := weight * math.Exp(-a * distance)
			result += distanceScore;
		}
    }
	return result;
}

func GetLiveSafeScore(c *gin.Context) {
	// 從請求中獲取 longitude 和 latitude 座標
	xStr := c.Query("x") // longitude
	yStr := c.Query("y") // latitude

	// 將座標轉換為 float64
	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": xStr})
		return
	}
	y, err := strconv.ParseFloat(yStr, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	result := GetWeightScore(x, y, "traffic_accident_location_view.geojson", 0.00284, 5)
    if (result > 100) {
		result = 100
	}
	// 返回 JSON
    c.JSON(200, gin.H{"score": result})
}