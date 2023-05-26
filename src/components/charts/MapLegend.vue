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
                <p><strong>{{ item.name }}</strong>{{ item.value ? `: ${item.value}` : '' }}</p>
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

    &-legend {
        display: grid;
        grid-template-columns: 1fr 1fr;
        width: 100%;

        &-item {
            display: flex;
            align-items: center;

            div {
                width: 1rem;
                margin-right: 10px;
            }

            p {
                font-size: 1rem;
            }
        }
    }
}
</style>