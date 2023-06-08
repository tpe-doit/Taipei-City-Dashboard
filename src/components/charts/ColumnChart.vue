<!-- Cleaned -->

<script setup>
import { ref } from 'vue';

const props = defineProps(['chart_config', 'activeChart', 'series'])

const chartOptions = ref({
    chart: {
        stacked: true,
        toolbar: {
            show: false
        },
    },
    colors: props.chart_config.color,
    dataLabels: {
        enabled: props.chart_config.categories ? false : true,
        offsetY: 20,
    },
    grid: {
        show: false,
    },
    legend: {
        show: props.chart_config.categories ? true : false,

    },
    plotOptions: {
        bar: {
            borderRadius: 5,
        },
    },
    stroke: {
        colors: ['#282a2c'],
        show: true,
        width: 2,
    },
    tooltip: {
        // The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.labels[dataPointIndex] + `${props.chart_config.categories ? '-' + w.globals.seriesNames[seriesIndex] : ''}` + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + ` ${props.chart_config.unit}` + '</span>' +
                '</div>'
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
        },
        type: 'category',
    },
})

</script>

<template>
    <div v-if="activeChart === 'ColumnChart'">
        <apexchart width="100%" height="270px" type="bar" :options="chartOptions" :series="series"></apexchart>
    </div>
</template>