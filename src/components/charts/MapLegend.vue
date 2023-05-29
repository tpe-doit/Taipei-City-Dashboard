<script setup>
const props = defineProps(['chart_config', 'activeChart', 'series'])
</script>

<template>
    <div class="maplegend">
        <div class="maplegend-legend">
            <div v-for="item in series" :key="item.name" class="maplegend-legend-item">
                <div
                    :style="{ backgroundColor: `${item.color}`, height: item.type === 'line' ? '0.4rem' : '1rem', borderRadius: item.type === 'circle' ? '50%' : '2px' }">
                </div>
                <div v-if="item.value">
                    <h5>{{ item.name }}</h5>
                    <h6>{{ item.value }} {{ chart_config.unit }}</h6>
                </div>
                <div v-else>
                    <h6>{{ item.name }}</h6>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
.maplegend {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    margin-top: -1rem;
    overflow: visible;

    &-legend {
        display: grid;
        grid-template-columns: 1fr 1fr;
        column-gap: 1rem;
        row-gap: 1rem;
        overflow: visible;
        width: 100%;

        &-item {
            display: flex;
            align-items: center;
            border: 1px solid var(--color-border);
            border-radius: 5px;
            padding: 5px 10px 5px 5px;
            box-shadow: 0px 0px 5px black;

            div:first-child {
                width: 1rem;
                margin-right: 0.75rem;
            }

            h5 {
                color: var(--color-complement-text);
                font-size: 0.75rem;
            }

            h6 {
                font-size: 1rem;
                font-weight: 400;
            }
        }
    }
}
</style>