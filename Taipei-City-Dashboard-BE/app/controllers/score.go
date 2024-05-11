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

const earthRadius = 6371 // 地球半徑 km公里

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad) * math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c
	return distance
}

// Point 代表一個經緯度的點
type Point struct {
	Latitude, Longitude float64
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
	x: x value for coordinates
	y: y value for coordinates
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
	// 從請求中獲取 x 和 y 座標
	xStr := c.Query("x")
	yStr := c.Query("y")

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
	result := GetWeightScore(x, y, "traffic_accident_location_view.geojson", 0.00084, 5)
    if (result > 100) {
		result = 100
	}
	// 返回 JSON
    c.JSON(200, gin.H{"score": result})
}