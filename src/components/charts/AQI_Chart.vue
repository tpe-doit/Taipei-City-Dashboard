<!-- Developed by Taipei Urban Intelligence Center 2023 -->

<script setup>
import { computed, ref } from 'vue';
import { useMapStore } from '../../store/mapStore';

const props = defineProps(['chart_config', 'activeChart', 'series', 'map_config']);

const mapStore = useMapStore();

const heatmapData = computed(() => {
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
			for (let i = 0; i < props.chart_config.categories.length; i++) {
				if (!output[props.chart_config.categories[i]]) {
					output[props.chart_config.categories[i]] = 0;
				}
				output[props.chart_config.categories[i]] += +serie.data[i];

				if (+serie.data[i] > highest) highest = +serie.data[i];
			}
		});
		sum = Object.values(output).reduce((partialSum, a) => partialSum + a, 0);
	}

	output.highest = highest;
	output.sum = sum;
	return output;
});

const colorScale = computed(() => {
	const ranges =[
		{
			from:0,
			to:50,
			name:"良好",
			color:"#00e400"
		},
		{
			from:51,
			to:100,
			color:"#D6D46D",
			name:"普通"
		},
		{
			from:101,
			to:150,
			color:"#ff7e00",
			name:"對敏感族群不健康"
		},
		{
			from:151,
			to:200,
			color:"#ff0000",
			name:"對所有族群不健康"
		},
		{
			from:201,
			to:300,
			color:"#8f3f97",
			name:"非常不健康"
		},
		{
			from:301,
			to:500,
			color:"#7e0023",
			name:"危害"
		}
	] 
	;
	return ranges;
});

const chartOptions = ref({
	chart: {
		stacked: true,
		toolbar: {
			show: false,
		},
	},
	dataLabels: {
		enabled: true,
		style: {
			fontSize: '12px',
			fontWeight: 'normal'
		}
	},
	grid: {
		show: false,
	},
	legend: {
		offsetY:15,
		showForNullSeries: false,
		show: true,
	},
	markers: {
		size: 3,
		strokeWidth: 0,
	},
	plotOptions: {
		heatmap: {
			enableShades: false,
			radius: 4,
			colorScale: {
				ranges: colorScale.value,
			}
		},
	},
	stroke: {
		show: true,
		width: 2,
		colors: ['#282a2c'],
	},
	tooltip: {
		custom: function ({ series, seriesIndex, dataPointIndex, w }) {
			// The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css
			return '<div class="chart-tooltip">' +
				'<h6>' + `${w.globals.seriesNames[seriesIndex]}-${w.globals.labels[dataPointIndex]}` + '</h6>' +
				'<span>' + `${series[seriesIndex][dataPointIndex]}` + `${props.chart_config.unit}` + '</span>' +
				'</div>';
		},
	},
	xaxis: {
		axisBorder: {
			show: false,
		},
		axisTicks: {
			show: false,
		},
		categories: props.chart_config.categories ? props.chart_config.categories : [],
		labels: {
			offsetY: 5,
			formatter: function (value) {
				return value.length > 7 ? value.slice(0, 6) + "..." : value;
			}
		},
		tooltip: {
			enabled: false
		},
		type: 'category',
	},
	yaxis: {
		max: function (max) {
			if (!props.chart_config.categories) {
				return max;
			}
			return heatmapData.value.highest;
		},
	}
});

const selectedIndex = ref(null);

function handleDataSelection(e, chartContext, config) {
	if (!props.chart_config.map_filter) {
		return;
	}
	let toFilter = `${config.dataPointIndex} -${config.seriesIndex}`;
	if (toFilter !== selectedIndex.value) {
		mapStore.addLayerFilter(`${props.map_config[0].index}-${props.map_config[0].type}`, props.chart_config.map_filter[0], props.chart_config.map_filter[1][config.dataPointIndex], props.map_config[0]);
		selectedIndex.value = toFilter;
	} else {
		mapStore.clearLayerFilter(`${props.map_config[0].index}-${props.map_config[0].type}`, props.map_config[0]);
		selectedIndex.value = null;
	}
}
</script>

<template>
	<div v-if="activeChart === 'AQI_Chart'" class="aqiChart">
		<div class="heatmapchart-title">
			<h5>年份</h5>
			<h6>{{props.chart_config.year}}</h6>
		</div>
		<apexchart width="100%" height="360px" type="heatmap" :options="chartOptions" :series="series"
			@dataPointSeplection="handleDataSelection"></apexchart>
	</div>
</template>
<style scoped lang="scss">
.heatmapchart {

	&-title {
		display: flex;
		justify-content: center;
		flex-direction: column;
		margin: -0.2rem 0 -1.5rem;

		h5 {
			color: var(--color-complement-text);
		}

		h6 {
			color: var(--color-complement-text);
			font-size: var(--font-m);
			font-weight: 400;
		}
	}
}
</style>