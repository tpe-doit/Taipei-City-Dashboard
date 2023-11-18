<!-- Developed by Taipei Urban Intelligence Center 2023 -->

<script setup>
import { ref, computed, onMounted , defineComponent , reactive   } from 'vue';
import { useMapStore } from '../../store/mapStore';

const props = defineProps(['chart_config', 'activeChart', 'series', 'map_config']);
const mapStore = useMapStore();

// color config : https://coolors.co/palette/f0ead2-dde5b6-adc178-a98467-6c584c
// https://coolors.co/palette/463f3a-8a817c-bcb8b1-f4f3ee-e0afa0
const chartOptions = ref({
	chart: {
		offsetY: 15,
		stacked: true,
		toolbar: {
			show: false
		},
	},
});

// add helper components

const dialogCardComponent = defineComponent({
	name: 'dialogCardComponent',
	props: {
		title: {
			type: String,
			required: true,
		},
		subtitle: {
			type: String,
			required: true,
		},
		unit:{
			type: String,
			default: props.chart_config.unit,
		}
	},
	template: `
		<div class="card">
			<div class="card-body">
				<h6 class="card-subtitle mb-2 text-muted">{{ title }}</h6>
				<h5 class="card-title center">{{ subtitle }}{{ unit }}</h5>
			</div>
		</div>
	`,
});

const svgPathComponent = defineComponent({
	name: 'svgPathComponent',
	props: {
		d: {
			type: String,
			default: '',
		},
		fill: {
			type: String,
			default: 'none',
		},
		onmouseenter: {
			type: Function,
			default: () => { console.log("onmouseenter"); },
		},
		onmouseleave: {
			type: Function,
			default: () => { },
		},
		onmousemove:{
			type: Function,
			default: () => { },
		},
	},
	template: `
		<path
			:d="d"
			:fill="fill"
			:onmouseenter="onmouseenter"
			:onmouseleave="onmouseleave"
			:onmousemove="onmousemove"
			fill-opacity="0.3"
		/>
	`,
});

const svgTextComponent = defineComponent({
	name: 'svgTextComponent',
	props: {
		x: {
			type: Number,
			default: 0,
		},
		y: {
			type: Number,
			default: 0,
		},
		fill: {
			type: String,
			default: 'none',
		},
		text: {
			type: String,
			default: '',
		},
		transform:{
			type: String,
			default: '',
		}
	},
	template: `
		<text
			:x="x"
			:y="y"
			:fill="fill"
			:font-size="fontSize"
			:font-family="fontFamily"
			:text-anchor="textAnchor"
			:alignment-baseline="middle"
			:transform="transform"
		>
			{{ text }}
		</text>
	`,
});


const svgRectComponent = defineComponent({
	name: 'svgRectComponent',
	props: {
		x: {
			type: Number,
			default: 0,
		},
		y: {
			type: Number,
			default: 0,
		},
		width: {
			type: Number,
			default: 0,
		},
		height: {
			type: Number,
			default: 0,
		},
		fill: {
			type: String,
			default: 'none',
		},
	},
	template: `
		<rect
			:x="x"
			:y="y"
			:width="width"
			:height="height"
			:fill="fill"
		/>
	`,
});


const chartHeight = computed(() => {
	return `${40 + props.series.length * 24}`;
});


let finalPathList = reactive([]);
let finalRectList = reactive([]);
let finalTextList = reactive([]);

const svgPathList = computed(() => {
  return finalPathList;
})

const svgRectList = computed(() => {
  return finalRectList;
})

const svgTextList = computed(() => {
  return finalTextList;
})


onMounted(() => {
	console.log("SankeyChart mounted");
	console.log(props.chart_config);
	console.log(props.series);	

	createSankey(props.series);

});


const generateSVGPath = (topRight, topLeft, bottomLeft, bottomRight) => {

	console.log("props.series.length" , props.series.length);
    // Constructing the SVG path string
    const pathData = [
        // Move to the starting point (top right)
        `M${topRight[0]} ${topRight[1]}`,
        // Draw a curve to the top left
        `Q${topLeft[0]} ${topLeft[1]},${topLeft[0]} ${topLeft[1]}`,
        // Draw a straight line to the bottom left
        `L${bottomLeft[0]} ${bottomLeft[1]}`,
        // Draw a curve to the bottom right
        `Q${bottomRight[0]} ${bottomRight[1]},${bottomRight[0]} ${bottomRight[1]}`,
        // Close the path
        'Z',
    ].join(' ');

    return pathData;
}

function createSankey(data){

	// data = props.series;

	// calculate DAG
	let nodes = [];
	let weightSumByNode = {}; // weightSum[from]=sum of weight of all edges from "from node"
	let weightSumByLayer = []; // weightSumByLayer[layer]=sum of weight of all edges in "layer"

	for (let i = 0; i < data.length; i++) {
		let from = data[i].from;
		let to = data[i].to;

		if (!weightSumByNode[from]) {
			weightSumByNode[from] = 0;
		}
		weightSumByNode[from] += data[i].weight;

		if (nodes.indexOf(from) == -1) {
			nodes.push(from);
		}
		if (nodes.indexOf(to) == -1) {
			nodes.push(to);
		}
	}


	// init Graph (adjacency list)
	let Graph = {};
	for(let i=0;i<nodes.length;i++){
		Graph[nodes[i]] = [];
	}

	for (let i = 0; i < data.length; i++) {
		let from = data[i].from;
		let to = data[i].to;
		let weight = data[i].weight;

		Graph[from].push({
			to: to,
			weight: weight
		});
	}

	console.log("Graph", Graph);

	// init inDegree
	let inDegree = {};
	for (let i = 0; i < nodes.length; i++) {
		inDegree[nodes[i]] = 0;
	}

	// calculate inDegree 
	for (let i = 0; i < data.length; i++) {
		let to = data[i].to;
		inDegree[to]++;
	}

	console.log("inDegree", inDegree);

	// run DAG
	// And add node list by layer
	let layerList = []; // type : [ [node1,node2,node3], [node4,node5] ]
	let layerCount = 0;

	let queue = [];
	for (let i = 0; i < nodes.length; i++) {
		if (inDegree[nodes[i]] == 0) {
			queue.push(nodes[i]);
		}
	}
	let layer = [];
	let dp = {};
	let step = 0;
	while (queue.length > 0) {
		dp[step++] = queue;
		layer = [];
		let size = queue.length;
		console.log("current que" , queue);

		for (let i = 0; i < size; i++) {
			let node = queue.shift();
			layer.push(node);
			if (Graph[node]) {
				console.log("Graph[node]" , Graph[node] );
				for (let key in Graph[node]) {
					let to = Graph[node][key].to;
					inDegree[to]--;
					if (inDegree[to] == 0) {
						queue.push(to);
					}
				}
			}
		}
		layerList.push(layer);
	}

	console.log("layerList" , dp );
	
	console.log("layerList" , layerList );
	layerCount = layerList.length;

	// calculate weightSumByLayer
	for (let i = 0; i < layerCount - 1; i++) {
		let layer = layerList[i];
		let sum = 0;
		for (let j = 0; j < layer.length; j++) {
			sum += weightSumByNode[layer[j]];
		}
		weightSumByLayer.push(sum);
	}
	// calculate last layer
	let lastLayerSum = 0;
	let lastLayerDict = {};

	for (let i = 0; i < layerList[layerCount - 1].length; i++) {
		console.log(`last layer` , layerList[layerCount - 1][i] );
		if (!lastLayerDict[layerList[layerCount - 1][i]]) {
			weightSumByNode[layerList[layerCount - 1][i]] = 0;
		}

		lastLayerDict[layerList[layerCount - 1][i]] = true;
	}

	console.log(`lastLayerDict:` , lastLayerDict )

	for (let i = 0; i < data.length; i++) {
		let to = data[i].to;
		console.log(`to:${to}`);
		if (lastLayerDict[to]) {
			console.log("last");
			lastLayerSum += data[i].weight;
			weightSumByNode[to] += data[i].weight;
		}
	}
	weightSumByLayer.push(lastLayerSum);

	// add lastLayer to weightSumByNode
	// for (let i = 0; i < layerList[layerCount - 1].length; i++) {
	//     let node = layerList[layerCount - 1][i];
	//     weightSumByNode[node] = lastLayerSum;
	// }


	// print layerList

	console.log("layerCount", layerCount);
	console.log("layerList", layerList);
	console.log("weightSumByNode", weightSumByNode);
	console.log("weightSumByLayer", weightSumByLayer);

	// total component height
	// let TOTAL_HEIGHT = svg.height.baseVal.value;
	// let TOTAL_WIDTH = svg.width.baseVal.value;
	let TOTAL_HEIGHT = parseInt(chartHeight.value)*0.4;
	let TOTAL_WIDTH = 320;


	// make the line curve


	let widthUnit = TOTAL_WIDTH / (layerCount * 2 - 1);
	let smallwidthUnit = widthUnit / 3;
	let edgeWidth = smallwidthUnit * 4;
	let nodeWidth = smallwidthUnit * 2;

	// let edgeW
	let paddingUnit = 20;

	let svgObjectsDict = {};
	let xPadding = smallwidthUnit / 2;

	// add vertical padding
	let maxLayerCount = 0;
	for (let i = 0; i < layerList.length; i++) {
		maxLayerCount = Math.max(maxLayerCount, layerList[i].length);
	}



	// calculate all node position & size
	for (let currentLayer = 0; currentLayer < layerCount; currentLayer++) {


		let lastY = 0;

		// add vertical padding
		let currentLayerCount = layerList[currentLayer].length;
		lastY += (maxLayerCount - currentLayerCount) * paddingUnit / 2;

		for (let i = 0; i < layerList[currentLayer].length; i++) {
			console.log("i", i);
			console.log("layerList[currentLayer]", layerList[currentLayer]);
			let node = layerList[currentLayer][i];

			console.log(`layer:${currentLayer} , node:${node}`);

			let rate = weightSumByNode[node] / weightSumByLayer[currentLayer];
			// let rate = weightSumByNodeSigmoid[node]/weightSumByLayerSigmoid[currentLayer];

			console.log(`rate:${rate}`);

			let height = TOTAL_HEIGHT * (rate);
			let width = nodeWidth;

			// let x = TOTAL_WIDTH*currentLayer/layerCount;
			let x = currentLayer * (edgeWidth + nodeWidth) + xPadding;

			// add padding
			lastY += (i == 0) ? 0 : paddingUnit;

			console.log(`country:${node} , rate:${rate} `);
			console.log(`x:${x} , y:${lastY} , width:${width} , height:${height}`);

			// draw rect & text later
			

			svgObjectsDict[node] = {
				x: x,
				y: lastY,
				layer: currentLayer,
				width: width,
				height: height,
				rate: rate,
				used_from_height: 0,
				used_to_height: 0,
			};


			lastY += height
		}
	}

	console.log("svgObjectsDict", svgObjectsDict);

	// draw curve path 

	for (let i = 0; i < data.length; i++) {
		let from = data[i].from;
		let to = data[i].to;
		let weight = data[i].weight;

		console.log("from", from);
		console.log("to", to);
		console.log("svgObjectsDict", svgObjectsDict);

		let x1 = svgObjectsDict[from].x + nodeWidth;
		let x2 = svgObjectsDict[to].x;
		

		// edge height calculation
		let from_rate = weight / weightSumByLayer[svgObjectsDict[from].layer];
		let from_height = TOTAL_HEIGHT * (from_rate);

		let to_rate = weight / weightSumByLayer[svgObjectsDict[to].layer];
		let to_height = TOTAL_HEIGHT * (to_rate);



		// draw curve polygon from x1,y1 to x2,y2 by `path` element


		let topLeft = [x1, svgObjectsDict[from].y + svgObjectsDict[from].used_to_height];
		let topRight = [x2, svgObjectsDict[to].y + svgObjectsDict[to].used_from_height];
		let bottomRight = [x2, svgObjectsDict[to].y + svgObjectsDict[to].used_from_height + to_height];
		let bottomLeft = [x1, svgObjectsDict[from].y + svgObjectsDict[from].used_to_height + from_height];

		let pathData = generateSVGPath(topRight, topLeft, bottomLeft, bottomRight);
		// const svgPathElement = document.createElementNS("http://www.w3.org/2000/svg", "path");
		// svgPathElement.setAttribute("d", pathData);
		// svgPathElement.setAttribute("fill", "gray");
		// if (data[i].negative) {
		// 	svgPathElement.setAttribute("fill", "red");
		// }

		// svgPathElement.setAttribute("fill-opacity", "0.2");

		// svg.appendChild(svgPathElement);

		let fillColor = props.chart_config.edge_color;
		if (data[i].negative) {
			fillColor = props.chart_config.negative_edge_color;
		}

		finalPathList.push({
			d: pathData,
			fill: fillColor,
			real_weight: data[i].real_weight,
			onmouseenter: (e) => {
				toolTipState.value.state = true;
				console.log("onmouseenter");
				console.log("toolTipState", toolTipState.value.state);

				console.log( e );

				// set mouse position to path dom element position
				console.log("e.target.getBoundingClientRect()", e.target.getBoundingClientRect());	
				mousePosition.value.x = e.target.getBoundingClientRect().x;
				mousePosition.value.y = e.target.getBoundingClientRect().y;

				// set dialog data
				targetDialog.value.title = `${from} -> ${to}`;
				targetDialog.value.real_weight = data[i].real_weight;
				targetDialog.value.negative = data[i].negative;

				// change hover color
				// let fillColor = props.chart_config.edge_hover_color;
				// if (data[i].negative) {
				// 	fillColor = props.chart_config.negative_edge_hover_color;
				// }
				e.target.setAttribute("fill", fillColor);
				e.target.setAttribute("fill-opacity", "0.9");


				tooltip.title = "testTitle";
				tooltip.subtitle = data[i].real_weight + " " + props.chart_config.unit;
				tooltip.items = [
					{
						name: "負向",
						value: data[i].negative + " " + props.chart_config.unit,
						rate: data[i].negative / data[i].real_weight * 100 + "%",
					},
					{
						name: "正向",
						value: data[i].positive + " " + props.chart_config.unit,
						rate: data[i].positive / data[i].real_weight * 100 + "%",
					},
				];
			},
			onmouseleave: (e) => {
				toolTipState.value.state = false;
				console.log("onmouseleave");
				console.log("toolTipState", toolTipState.value.state);
				console.log("e", e);

				// change hover color back
				// let fillColor = props.chart_config.edge_color;
				// if (data[i].negative) {
				// 	fillColor = props.chart_config.negative_edge_color;
				// }
				e.target.setAttribute("fill", fillColor);
				e.target.setAttribute("fill-opacity", "0.3");
			},
			onmousemove: (e) => {
				// console.log("onmousemove");
				// console.log("e", e);

				// mousePosition.value.x = e.pageX;
				// mousePosition.value.y = e.pageY;
				
			}
		});

		console.log("real_weight", data[i].real_weight);




		svgObjectsDict[from].used_to_height += from_height;
		svgObjectsDict[to].used_from_height += to_height;

		console.log(`from:${from} , to:${to} , weight:${weight} , i: ${i}`);

		console.log("topLeft", topLeft);
		console.log("topRight", topRight);
		console.log("bottomRight", bottomRight);
		console.log("bottomLeft", bottomLeft);
	}

	// draw rect & text
	// loop through svgObjectsDict
	for (let key in svgObjectsDict) {
		let x = svgObjectsDict[key].x;
		let y = svgObjectsDict[key].y;
		let width = svgObjectsDict[key].width;
		let height = svgObjectsDict[key].height;
		let rate = svgObjectsDict[key].rate;


		finalRectList.push({
			x: x,
			y: y,
			width: width,
			height: height,
			fill: props.chart_config.layer_colors[svgObjectsDict[key].layer],
		});

		// add text to rect

		finalTextList.push({
			x: x,
			y: y,
			fill: "white",
			fontSize: "12px",
			fontFamily: "sans-serif",
			textAnchor: "middle",
			text: key,
			transform: `translate(${width / 2},${height / 2})`,
		});
	}
}
	
	// handle tooltip
	let tooltip = ref({
		title: '',
		subtitle: '',
		items: [],
	});
	const toolTipState = ref({state: false});
	const showTooltip = computed(() => {
		console.log("showTooltip", toolTipState);
		return toolTipState.value.state;
	});

	const targetDialog = ref({
		title: '',
		real_weight: '',
		negative: false,
	});
	// const DialogColor = ref(props.chart_config.color[0]);
	const mousePosition = ref({ x: null, y: null });
	const selectedIndex = ref(null);

	// Parse Dialog Data (to support 2D or 3D data)
	const DialogData = computed(() => {
		let output = {};
		let highest = 0;
		let sum = 0;
		if (props.series.length === 1) {
			props.series[0].data.forEach((item) => {
				output[item.x] = item.y;
				if (item.y > highest) {
					highest = item.y;
				}
				sum += item.y;
			});
		} else {
			props.series.forEach((serie) => {
				for (let i = 0; i < 12; i++) {
					if (!output[props.chart_config.categories[i]]) {
						output[props.chart_config.categories[i]] = 0;
					}
					output[props.chart_config.categories[i]] += +serie.data[i];
				}
			});
			highest = Object.values(output).sort(function (a, b) { return b - a; })[0];
			sum = Object.values(output).reduce((partialSum, a) => partialSum + a, 0);
		}

		output.highest = highest;
		output.sum = sum;
		return output;
	});
	const tooltipPosition = computed(() => {
		return { 'left': `${mousePosition.value.x - 54}px`, 'top': `${mousePosition.value.y - 54}px` };
		// return { 'left': `${mousePosition.value.x}px`, 'top': `${mousePosition.value.y}px` };
	});

	function handleDataSelection(index) {
		if (!props.chart_config.map_filter) {
			return;
		}
		if (index !== selectedIndex.value) {
			mapStore.addLayerFilter(`${props.map_config[0].index}-${props.map_config[0].type}`, props.chart_config.map_filter[0], props.chart_config.map_filter[1][index]);
			selectedIndex.value = index;
		} else {
			mapStore.clearLayerFilter(`${props.map_config[0].index}-${props.map_config[0].type}`);
			selectedIndex.value = null;
		}
	}

</script>

<template>
	<div v-if="activeChart === 'SankeyChart'">
		 <div width="100%" :height="chartHeight" :options="chartOptions" :series="series"
			@dataPointSelection="handleDataSelection"> 

			<svg id="sankey" width="100%" :height="chartHeight" xmlns="http://www.w3.org/2000/svg">

				<svgPathComponent v-for="path in svgPathList" :d="path.d" :fill="path.fill" :stroke="path.stroke"
					:stroke-width="path.strokeWidth" :onmouseenter="path.onmouseenter" :onmouseleave="path.onmouseleave"
					:onmousemove="path.onmousemove" />
				<svgRectComponent v-for="rect in svgRectList" :x="rect.x" :y="rect.y" :width="rect.width"
					:height="rect.height" :fill="rect.fill" />
				<svgTextComponent v-for="text in svgTextList" :x="text.x" :y="text.y" :fill="text.fill"
					:font-size="text.fontSize" :font-family="text.fontFamily" :text-anchor="text.textAnchor"
					:alignment-baseline="text.alignmentBaseline" :transform="text.transform" :text="text.text" />
				
				
			</svg>
			<div v-if="showTooltip" >
				<h1>測試 ouo</h1>
			</div>
			<div v-if="1" >
				<h1>測試 aaa</h1>
			</div>

			<dialogCardComponent v-if="showTooltip" :title="targetDialog.title" :subtitle="targetDialog.real_weight" :style="tooltipPosition">
			</dialogCardComponent>

		</div>
	</div>
</template>

<style>

#sankey {
	overflow: visible !important;
	z-index: 999;
}

/* dialog card style */

.card {
	position: fixed;
	border: none;
	border-radius: 10px;
	background-color: #485159;
	color: white;
	box-shadow: rgba(50, 50, 93, 0.25) 0px 50px 100px -20px, rgba(0, 0, 0, 0.3) 0px 30px 60px -30px, rgba(10, 37, 64, 0.35) 0px -2px 6px 0px inset;
	z-index: 99999;
}

.card-body {
	padding: 0.5rem 1rem;
}

.card-title {
	font-size: 1rem;
}

.card-subtitle {
	font-size: 0.8rem;
	color: #dadde3;
}


</style>