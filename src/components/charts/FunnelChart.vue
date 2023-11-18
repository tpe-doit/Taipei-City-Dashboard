<script setup>
import { ref } from "vue";

const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);

const chartOptions = ref({
	chart: {
		type: "bar",
		height: 300,
		toolbar: false,
	},
	plotOptions: {
		bar: {
			borderRadius: 0,
			horizontal: true,
			isFunnel: true,
		},
	},
	dataLabels: {
		enabled: true,
		formatter: function (val, opt) {
			return opt.w.globals.labels[opt.dataPointIndex] + ":  " + val;
		},
		dropShadow: {
			enabled: true,
		},
	},
	xaxis: {
		categories: props.chart_config.categories,
	},
	legend: {
		show: false,
	},
	tooltip: {
		custom: function ({ series, seriesIndex, dataPointIndex, w }) {
			const value = props.chart_config.useAbs
				? Math.abs(series[seriesIndex][dataPointIndex])
				: series[seriesIndex][dataPointIndex];
			return (
				'<div class="chart-tooltip">' +
				"<h6>" +
				w.globals.labels[dataPointIndex] +
				"</h6>" +
				"<span>" +
				value +
				` ${props.chart_config.unit}` +
				"</span>" +
				"</div>"
			);
		},
		followCursor: true,
	},
});
</script>

<template>
	<div v-if="activeChart === 'FunnelChart'">
		<apexchart
			width="100%"
			type="bar"
			:options="chartOptions"
			:series="series"
		></apexchart>
	</div>
</template>
