<script setup>
import { ref, computed } from 'vue'
import { useMapStore } from '../../store/mapStore';

const props = defineProps(['chart_config', 'activeChart', 'series', 'map_config']);
const mapStore = useMapStore();

const chartOptions = ref({
	chart: {
		// offsetY: 15,
		stacked: true,
		toolbar: {
			show: false,
		}
	},
	colors: props.chart_config.color,
	dataLabels: {
		offsetX: 20,
		textAnchor: 'start',
	},
	grid: {
		show: false,
	},
	legend: {
		show: false,
	},
	plotOptions: {},
	stroke: {
		colors: ['#282a2c'],
		show: true,
		width: 0,
	},
	tooltip: {
		custom: function ({ series, seriesIndex, dataPointIndex, w }) {
			return '<div class="chart-tooltip">' +
				'<h6>' + w.globals.labels[dataPointIndex] + '</h6>' +
				'<span>' + series[seriesIndex][dataPointIndex] + ` ${props.chart_config.unit}` + '</span>' +
				'</div>';
		},
		followCursor: true,
	},
	xaxis: {
		axisBorder: {
			show: false,
		},
		axisTicks: {
			show: false,
		},
		labels: {
			show: false,
		},
		type: 'category',
	},
	yaxis: {
		labels: {
			formatter: function (value) {
				return value.length > 7 ? value.slice(0, 6) + "..." : value;
			},
		},
	},
})
</script>

<template>
	<div v-if="activeChart === 'RangeAreaChart'">
		<apexchart width="100%" type="rangeArea" :options="chartOptions" :series="series">
		</apexchart>
	</div>
</template>