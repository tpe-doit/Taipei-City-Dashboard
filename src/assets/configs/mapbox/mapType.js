export const mapBasicFormObj = {
    index: '',
    symbol: '',
    property: [],
    title: '',
    cover: false,
    color: '',
    paint: {},
}

const mapSource = [
    {
        type: 'raster',
        title: 'Geoserver',
        mapFormTemplate: {
            ...mapBasicFormObj
        }
    },
    {
        type: 'geoJson',
        title: 'API',
        mapFormTemplate: {
            ...mapBasicFormObj,
            url: '',
            form_data: {}
        }
    }
]

/**
 * 花博 expopark: /public/datas/expopark.geojson
 * [{"geoJson":{"color":"#1b71b0","index":"event","property":[{"key":"location","name":"展區位置"},{"key":"event","name":"活動名稱"}],"symbol":"line","title":"展區位置"}}]
 * 
 * 士林 shilin: /public/datas/expopark.geojson
 */
export const MapPositions = {
    center: [121.536609, 25.044808],
    taipeiBounds: [
        [121.3870596781498, 24.95733863075891], // Southwest coordinates
        [121.6998231749096, 25.21179993640203] // Northeast coordinates
    ],
    defaultZoom: 12.5,  // Less than 15 GetFeatureInfo does not work
    defaultMinZoom: 11,
    defaultMaxZoom: 22,
    CityHall: [121.56719891052944, 25.036536501522235],
    Eastern: [121.55676686461811, 25.038745980403505],
    Dadaocheng: [121.51013649269608, 25.055775323336594],
    Yangmingshan: [121.55180953513364, 25.12898719167495],
    ExpoPark: [121.5228860863223,25.068515771454074],
    Shilin: [121.52231934552663,25.09131257257215],
    n: [121.53799903832744, 25.092301189452208],
    e:[ 121.5858901814368, 25.046019127628448],
    w:[121.48719722062226, 25.050157512489278],
    s:[121.52908519223917, 25.005905715593244]
}

export const rainFallLayerKey = [78.8, 100, 130]
export const rainFallLayers = {
    '78.8':{
        "index": "patrol_flood_78_8",
        "title":"近一小時降雨78.8mm可能淹水範圍",
    },
    '100':{
        "index": "patrol_flood_100",
        "title":"近一小時降雨100mm可能淹水範圍",
    },
    '130':{
        "index": "patrol_flood_130",
        "title":"近一小時降雨130mm可能淹水範圍",
    }
}
export const rainFallConfig = {
    "type": 'fill',
    "layout": {'visibility': 'none'},
    "paint": {
        "fill-color":[
            "match",["get","name"],
            "1.0 m ~ 2.0 m","#352ad1",
            "0.30 m ~ 1.0 m","#009ff4",
            "0.15 m ~ 0.30 m","#baf4f5",
            "#352ad1"
        ],
        "fill-opacity": [
            "interpolate",["linear"],["zoom"],
            11,0.5,
            13.5,0.3
        ]
    }
}
export const weatherLayers = {
    'temperature':{
        "index": "temperature",
        "title":"溫度",
        "color": [
            "interpolate",["linear"],["get","temperature"],
            -20,"#b1d8ea",
            20,"#ffc107",
            40,"#f44336"
        ],
        "height": 1
    },
    'humidity':{
        "index": "humidity",
        "title":"濕度",
        "color": [
            "interpolate",["linear"],["get","humidity"],
            0,"#ddd",
            1,"#2196f3"
        ],
        "height": 1
    },
    'rainfall_24hour':{
        "index": "rainfall_24hour",
        "title":"日累積雨量",
        "color": [
            "interpolate",["linear"],["get","rainfall_24hour"],
            0,"#ddd",
            1,"#2196f3"
        ],
        "height": [
            "interpolate",["linear"],["get","rainfall_24hour"],
            0,1,
            200,200
        ]
    }
}

export const symbolList = [
    {
        label: '點',
        options: [
            {value:'circle',name: '圓點'},
            {value:'circleMiddle',name: '圓點(中)'},
            {value:'circleBig',name: '圓點(大)'},
            {value:'heatmap',name: '熱力圖＆圓點'}
        ]
    },{
        label: '線',
        options: [
            {value:'line',name: '實線(粗)'},
            {value:'lineThin',name: '實線(細)'},
            {value:'dasharray',name: '虛線'}
        ]
    },{
        label: '面',
        options: [
            {value:'fill',name: '平面'},
            {value:'fill-extrusion',name: '立體'}
        ]
    },{
        label: 'Icon',
        options: [
            {value:'cross1',name: '交叉'},
            {value:'triangle1',name: '三角形(白)'},
            {value:'triangle2',name: '三角形(綠)'},
            {value:'bus',name: '公車'},
            {value:'metro',name: '捷運'},
            {value:'bike',name: '腳踏車'},
            {value:'weather',name: '氣象'},
            {value:'symbol',name: '其他設定'}
        ]
    }
]

export default mapSource