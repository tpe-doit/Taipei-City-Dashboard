<!-- Developed by Taipei Urban Intelligence Center 2023 -->

<script setup>
import { ref } from "vue";
import { useMapStore } from "../../store/mapStore";

const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);
const mapStore = useMapStore();
function generateData(baseval, count, yrange) {
	var i = 0;
	var series = [];
	while (i < count) {
		var x = Math.floor(Math.random() * (750 - 1 + 1)) + 1;
		var y =
			Math.floor(Math.random() * (yrange.max - yrange.min + 1)) +
			yrange.min;
		var z = Math.floor(Math.random() * (75 - 15 + 1)) + 15;

		series.push([x, y, z]);
		baseval += 86400000;
		i++;
	}
	return series;
}

const chartOptions = {
	series: [
		{
			name: "Bubble1",
			data: generateData(new Date("11 Feb 2017 GMT").getTime(), 20, {
				min: 10,
				max: 60,
			}),
		},
		{
			name: "Bubble2",
			data: generateData(new Date("11 Feb 2017 GMT").getTime(), 20, {
				min: 10,
				max: 60,
			}),
		},
		{
			name: "Bubble3",
			data: generateData(new Date("11 Feb 2017 GMT").getTime(), 20, {
				min: 10,
				max: 60,
			}),
		},
		{
			name: "Bubble4",
			data: generateData(new Date("11 Feb 2017 GMT").getTime(), 20, {
				min: 10,
				max: 60,
			}),
		},
	],
	chart: {
		height: 350,
		type: "bubble",
		toolbar: false,
	},
	dataLabels: {
		enabled: false,
	},
	fill: {
		opacity: 0.8,
	},
	title: {
		text: "Simple Bubble Chart",
	},
	xaxis: {
		tickAmount: 12,
		type: "category",
	},
	yaxis: {
		max: 70,
	},
};

const selectedIndex = ref(null);

function handleDataSelection(e, chartContext, config) {
	if (!props.chart_config.map_filter) {
		return;
	}
	if (config.dataPointIndex !== selectedIndex.value) {
		mapStore.addLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`,
			props.chart_config.map_filter[0],
			props.chart_config.map_filter[1][config.dataPointIndex]
		);
		selectedIndex.value = config.dataPointIndex;
	} else {
		mapStore.clearLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`
		);
		selectedIndex.value = null;
	}
}
</script>

<template>
	<div v-if="activeChart === 'BubbleChart'">
		<apexchart
			width="100%"
			type="bubble"
			height="300px"
			:options="chartOptions"
			@dataPointSelection="handleDataSelection"
			:series="series"
		></apexchart>
	</div>
</template>
