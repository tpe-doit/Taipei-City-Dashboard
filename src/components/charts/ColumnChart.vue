<script setup>
import { defineProps, ref, computed } from 'vue';
const props = defineProps(['chart_config', 'activeChart', 'series'])

const options = ref({
    legend: {
        show: props.chart_config.categories ? true : false,

    },
    stroke: {
        show: true,
        colors: ['#282a2c'],
        width: 2,
    },
    tooltip: {
        custom: function ({ series, seriesIndex, dataPointIndex, w }) {
            return '<div class="chart-tooltip">' +
                '<h6>' + w.globals.labels[dataPointIndex] + `${props.chart_config.categories ? '-' + w.globals.seriesNames[seriesIndex] : ''}` + '</h6>' +
                '<span>' + series[seriesIndex][dataPointIndex] + ` ${props.chart_config.unit}` + '</span>' +
                '</div>'
        }
    },
    dataLabels: {
        enabled: props.chart_config.categories ? false : true,
        offsetY: 20
    },
    plotOptions: {
        bar: {
            borderRadius: 2,
        }
    },
    colors: props.chart_config.color,
    grid: {
        show: false,
    },
    chart: {
        toolbar: {
            show: false
        },
        stacked: true,

    },
    xaxis: {
        type: 'category',
        categories: props.chart_config.categories ? props.chart_config.categories : [],
        axisTicks: {
            show: false,
        },
        labels: {
            offsetY: 5,
        }
    },
})

</script>

<template>
    <div v-if="activeChart === 'ColumnChart'">
        <apexchart width="100%" height="270px" type="bar" :options="options" :series="series"></apexchart>
    </div>
</template>

<style scoped></style>