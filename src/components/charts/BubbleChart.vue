<!-- Developed by Taipei Urban Intelligence Center 2023 -->

<script setup>
const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);

const chartOptions = {
	chart: {
		type: "bubble",
		toolbar: false,
	},
	dataLabels: {
		enabled: false,
	},
	fill: {
		opacity: 0.8,
	},
	xaxis: {
		tickAmount: 12,
		type: "category",
	},
	yaxis: {
		max: props.chart_config.yaxis,
	},
	tooltip: {
		// The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css
		custom: function ({ series, seriesIndex, dataPointIndex, w }) {
			return (
				'<div class="chart-tooltip">' +
				"<h6>" +
				series[seriesIndex][dataPointIndex] +
				"</h6>" +
				"<span>" +
				w.globals.labels[dataPointIndex] +
				` ${props.chart_config.unit}` +
				"</span>" +
				"</div>"
			);
		},
		followCursor: true,
	},
};
</script>

<template>
	<div v-if="activeChart === 'BubbleChart'">
		<apexchart
			width="100%"
			type="bubble"
			height="300px"
			:options="chartOptions"
			:series="series"
			:colors="chartOptions.colors"
		></apexchart>
	</div>
</template>
