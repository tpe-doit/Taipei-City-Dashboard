<!-- Developed by Taipei Urban Intelligence Center 2023 -->

<script setup>
import { ref } from 'vue';
import { useMapStore } from '../../store/mapStore';
import dayjs from 'dayjs';

const props = defineProps(['chart_config', 'activeChart', 'series', 'map_config']);
const mapStore = useMapStore();

const chartOptions = ref({
	chart: {
		toolbar: {
			show: false,
			tools: {
				zoom: false,
			},
		},
	},
	colors: props.chart_config.color,
	dataLabels: {
		enabled: false,
	},
	grid: {
		show: false,
	},
	legend: {
		show: props.series.length > 1 ? true : false,
		onItemClick: {
			toggleDataSeries: false,
		},
	},
	markers: {
		hover: {
			size: 4,
		},
		size: 3,
		strokeWidth: 0,
	},
	stroke: {
		colors: props.chart_config.color,
		curve: 'smooth',
		show: true,
		width: 2,
	},
	tooltip: {
		shared: false,
		// intersect: true,
		custom: function ({ series, seriesIndex, dataPointIndex, w }) {
			// The class 'chart-tooltip' could be edited in /assets/styles/chartStyles.css
			return (
				'<div class="chart-tooltip">' +
				"<h6>" +
				parseTime(w.config.series[seriesIndex].data[dataPointIndex].x, props.chart_config.tooltipTimeFormat) +
				' - ' + w.globals.seriesNames[seriesIndex] +
				"</h6>" +
				"<span>" +
				`${series[seriesIndex][dataPointIndex]} ${props.chart_config.unit}` +
				"</span>" +
				"</div>"
			);
		},
	},
	xaxis: {
		axisBorder: {
			color: "#555",
			height: "0.8",
		},
		axisTicks: {
			show: false,
		},
		crosshairs: {
			show: false,
		},
		tooltip: {
			enabled: false,
		},
		type: "datetime",
	},
});

function parseTime(time, format) {
	return dayjs(time).format(format);
}

const selectedIndex = ref(null);

function handleDataSelection(e, chartContext, config) {
	if (!props.chart_config.map_filter)
		return;

	if (e.target.className !== 'apexcharts-legend-text') {
		const layer_id = `${props.map_config[0].index}-${props.map_config[0].type}`;
		if (selectedIndex.value != null) {
			selectedIndex.value = null;
			mapStore.restorePaintProperty(layer_id, props.map_config[0], 'color');
		}
		return;
	}

	const dataPointIndex = e.target.getAttribute('i');
	const layer_id = `${props.map_config[0].index}-${props.map_config[0].type}`;
	if (dataPointIndex !== selectedIndex.value) {
		selectedIndex.value = dataPointIndex;
		const prop = props.chart_config.map_filter[1][dataPointIndex];
		if(prop==null)
			return;
		mapStore.setPaintProperty(layer_id, props.map_config[0], 'color', '#505050', [props.chart_config.map_filter[0], prop], true);
	}
}
</script>

<template>
	<div v-if="activeChart === 'TimelineSeparateChartCustom'">
		<apexchart
			width="100%"
			height="260px"
			type="line"
			:options="chartOptions"
			:series="series"
			@mouseMove="handleDataSelection"
			@mouseLeave="handleDataSelection"
		></apexchart>
	</div>
</template>
