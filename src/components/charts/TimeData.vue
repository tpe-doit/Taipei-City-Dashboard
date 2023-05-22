<script setup>
import { onMounted, ref } from 'vue';
import { parseTimeData } from '../../assets/utilityFunctions/parseChartData'
const props = defineProps({
    content: Object, defaultType: {
        type: String,
        default: 'separate'
    }
})

const chartOptions = ref({
    chart: {
        plotBackgroundColor: null,
        plotBorderWidth: null,
        backgroundColor: null,
        plotShadow: false,
        spaceTop: 0,
        height: '70%',
    },
    credits: {
        enabled: false,
    },
    exporting: {
        enabled: false,
    },
    legend: {
        enabled: true,
        itemStyle: {
            color: 'white'
        },
        itemHoverStyle: {
            color: '#888787'
        },
        maxHeight: 64,
        navigation: {
            arrowSize: 4
        }
    },
    plotOptions: {
        series: {
            stacking: undefined,
            marker: {
                symbol: 'circle'
            },
        },
        line: {
            allowPointSelect: true,
            cursor: "pointer",
            dataLabels: {
                enabled: true,
                format: "{point.y}",
            },
            borderWidth: 2,
            borderColor: "#282a2c",
        },
        areaspline: {
            allowPointSelect: true,
            cursor: "pointer",
            dataLabels: {
                enabled: true,
                format: "{point.y}",
            },
            borderWidth: 2,
            borderColor: "#282a2c",
        }
    },
    title: {
        text: undefined,
    },
    series: [],
    xAxis: {
        text: null,
        type: 'category',
        labels: {
            formatter: function () {
                if (this.value.length > 6) {
                    return this.value.substring(0, 6) + '...';
                }
                return this.value;
            },
            style: {
                color: 'white',
            }
        }
    },
    yAxis: {
        title: undefined,
        labels: {
            enabled: false
        },
        gridLineWidth: 0,
    },
    tooltip: {
        pointFormat: "<b>{series.name}</b>:{point.y}",
        backgroundColor: "#090909",
        style: {
            color: "#888787",
        },
    },
})

onMounted(() => {
    const dataToBeParsed = props.content.index === '8719a9b0' ? props.content.chartData[1].data : props.content.chartData[0].data
    const [newData, time] = parseTimeData(dataToBeParsed, props.content.request_list[0].form_data)
    chartOptions.value.series = newData

    if (props.content.request_list[0].color && props.content.request_list[0].color.length > 1) {
        chartOptions.value.colors = props.content.request_list[0].color
    }
    if (chartOptions.value.series.length <= 1 || chartOptions.value.series.length >= 4) {
        chartOptions.value.legend.enabled = false
    }

    if (props.defaultType === 'separate') {
        toSeparate();
    } else if (props.defaultType === 'combined') {
        toCombined();
    }
})

function toSeparate() {
    chartOptions.value.chart.type = 'line'
    chartOptions.value.plotOptions.series.stacking = undefined
}
function toCombined() {
    chartOptions.value.chart.type = 'areaspline'
    chartOptions.value.plotOptions.series.stacking = 'normal'
}

</script>

<template>
    <div class="timedata">
        <highcharts :options="chartOptions" :style="{ height: '90%', marginTop: '1.5rem' }"></highcharts>
        <div class="timedata-control">
            <button @click="toSeparate">比較</button>
            <button @click="toCombined">堆疊</button>
        </div>
    </div>
</template>

<style scoped lang="scss">
.timedata {
    min-height: 100%;
    max-height: 100%;
    overflow-y: scroll;

    &-control {
        position: absolute;
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        position: absolute;
        top: 0;

        button {
            background-color: rgb(77, 77, 77);
            padding: 4px 4px;
            border-radius: 5px;
            transition: color 0.2s, opacity 0.2s;
            font-size: var(--font-s);
            margin: 0 4px;
            color: var(--color-complement-text);
            opacity: 0.25;
            text-align: center;

            &:hover {
                color: white;
                opacity: 1;
            }
        }
    }

}
</style>