<script setup>
const props = defineProps(['direction', 'weight'])
const numToColor = {
    "1": "#acc22a",
    "2": "#fcd100",
    "3": "#ed8502",
    "4": "#ce161a",
}

</script>

<template>
    <div class="metrocardensity">
        <div v-if="weight" class="metrocardensity-container">
            <div class="metrocardensity-item" v-for="(item, index) in weight.y"
                :style="{ backgroundColor: numToColor[item] }">{{
                    index + 1
                }}</div>
        </div>
        <span v-else>{{ direction }}</span>
    </div>
</template>

<style scoped lang="scss">
.metrocardensity {
    display: flex;
    align-items: center;
    justify-content: center;

    span {
        font-family: var(--font-icon);
        color: var(--color-complement-text);
        pointer-events: none;
        user-select: none;
    }

    &-container {
        display: grid;
        grid-template-columns: 1fr 1fr 1fr 1fr 1fr 1fr;
        column-gap: 2px;
    }

    &-item {
        font-size: 0.8rem;
        width: 1rem;
        height: 1rem;
        border-radius: 2px;
        text-align: center;
        color: black;
        pointer-events: none;
        user-select: none;
    }
}
</style>