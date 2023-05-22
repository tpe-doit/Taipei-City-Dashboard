import { hexToHSL } from "../../utilityFunctions/colorConvert.js";
import * as mapLayerStyle from "./mapConfig.js";

export const ParseMapLayer = (MapConfig, MapLabel = {}) => {
  if (!MapConfig) return;
  const MapLayerIndex = MapConfig.index;
  const mapIconColor =
    MapLabel && MapLabel.color
      ? MapLabel.color
      : MapConfig.color
      ? MapConfig.color
      : "#ddd";

  let symbolType =
    MapLabel && MapLabel.symbol ? MapLabel.symbol : MapConfig.symbol;
  let symbolPaint = {};
  let symbolLayout = {};
  let mapExtraLayer = {};
  let loadImage = "";

  switch (symbolType) {
    case "circle":
      symbolPaint = {
        ...mapLayerStyle.circleCommonStyle,
        "circle-color": mapIconColor,
      };
      break;

    case "circle_stroke":
      symbolType = "circle";
      symbolPaint = {
        ...mapLayerStyle.circleCommonStyle,
        "circle-color": mapIconColor,
      };
      break;

    case "heatmap":
      symbolPaint = {
        ...mapLayerStyle.circleHeatmapStyle,
        "circle-color": mapIconColor,
      };
      break;

    case "line":
      symbolPaint = {
        ...mapLayerStyle.lineWideStyle,
        "line-color": mapIconColor,
      };
      break;

    case "lineThin":
      symbolPaint = {
        ...mapLayerStyle.lineSuperWideStyle,
        "line-color": mapIconColor,
      };
      break;

    case "dasharray":
      symbolPaint = {
        ...mapLayerStyle.lineWideStyle,
        "line-color": mapIconColor,
        "line-dasharray": [2, 2],
      };
      break;

    case "fill":
      symbolPaint = {
        ...mapLayerStyle.fillCommonStyle,
        "fill-color": mapIconColor,
      };
      break;

    case "fill-extrusion":
      break;

    default:
      if (symbolType === "triangle_white" || symbolType === "triangle_green") {
        symbolLayout = {
          ...mapLayerStyle.symbolShapeSmallStyle,
          "icon-image": symbolType,
        };
      } else if (symbolType === "cross_normal" || symbolType === "cross_bold") {
        symbolLayout = {
          ...mapLayerStyle.symbolIconSmallStyle,
          "icon-image": symbolType,
        };
      } else if (symbolType === "bus") {
        symbolLayout = mapLayerStyle.symbolBusStyle;
      } else if (symbolType === "bike") {
        symbolLayout = {
          ...mapLayerStyle.symbolTextHaloStyle.layout,
          "icon-image": symbolType,
        };
        symbolPaint = { ...mapLayerStyle.symbolTextHaloStyle.paint };
      } else if (symbolType === "metro") {
        symbolLayout = {
          ...mapLayerStyle.symbolTextHaloStyle.layout,
          "icon-image": symbolType,
        };
        symbolPaint = { ...mapLayerStyle.symbolTextHaloStyle.paint };
      } else if (symbolType === "symbol") {
        //     symbolLayout = { ...mapLayerStyle.symbolHaloStyle.layout, 'icon-image': symbolType + '-15'}
        //     symbolPaint = { ...mapLayerStyle.symbolHaloStyle.paint, "text-halo-color": mapIconColor }
      }
      loadImage = symbolType;
      symbolType = "symbol";
      break;
  }
  symbolLayout = {
    ...symbolLayout,
    visibility: "visible",
  };

  if (MapConfig.layout) {
    symbolLayout = {
      ...symbolLayout,
      ...MapConfig.layout,
    };
  }
  if (MapConfig.paint) {
    symbolPaint = {
      ...symbolPaint,
      ...MapConfig.paint,
    };
  }

  const mapLayerConfig = {
    id: MapLayerIndex,
    source: `${MapLayerIndex}_source`,
    paint: symbolPaint,
    layout: symbolLayout,
  };

  // if(MapConfig.filter && MapConfig.filter.target === 'hr'){
  //     const refer = dayjs().hour()
  //     mapLayerConfig.filter = ['==', ['get', 'hr'], refer]
  // }

  // Add extra layers
  if (symbolType === "heatmap") {
    symbolType = "circle";
    //https://docs.mapbox.com/mapbox-gl-js/example/heatmap-layer/
    const heatmap_hsl = hexToHSL(mapIconColor);
    mapExtraLayer = {
      ...mapLayerConfig,
      id: `heatmap_${MapLayerIndex}`,
      type: "heatmap",
      paint: {
        ...mapLayerStyle.heatmapCommonStyle,
        "heatmap-intensity": [
          "interpolate",
          ["linear"],
          ["zoom"],
          10.99,
          0,
          11,
          0.01,
          15,
          0.2,
        ],
        "heatmap-radius": 10,
        "heatmap-color": [
          "step",
          ["heatmap-density"],
          `hsla(${heatmap_hsl}, 0%)`,
          0.1,
          `hsla(${heatmap_hsl}, 25%)`,
          0.6,
          `hsla(${heatmap_hsl}, 50%)`,
          1,
          `hsla(${heatmap_hsl}, 75%)`,
        ],
      },
      interactive:
        MapConfig.interactive && MapConfig.interactive.affected
          ? MapConfig.interactive.affected
          : null,
    };
  } else if (symbolType === "line" || symbolType === "dasharray") {
    symbolType = "line";
    mapExtraLayer = {
      ...mapLayerConfig,
      id: `linefill_${MapLayerIndex}`,
      type: "fill",
      paint: {
        ...mapLayerStyle.fillCommonStyle,
        "fill-color": `hsla(${hexToHSL(mapIconColor)}, 50%)`,
      },
      interactive:
        MapConfig.interactive && MapConfig.interactive.affected
          ? MapConfig.interactive.affected
          : null,
    };
  } else if (symbolType === "lineThin") {
    symbolType = "line";
  }
  return {
    main: {
      ...mapLayerConfig,
      type: symbolType,
    },
    extra: mapExtraLayer,
    loadImage: loadImage,
    interactive:
      MapConfig.interactive && MapConfig.interactive.affected
        ? MapConfig.interactive.affected
        : null,
  };
};
