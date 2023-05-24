<script setup>
const props = defineProps({
    icon: String, text: String, mode: {
        type: String,
        default: 'outline'
    }
})
</script>

<template>
    <div :class="{ 'componenttag': true, 'componenttag-fill': mode === 'fill' }">
        <span>{{ icon }}</span>
        <p>{{ text }}</p>
    </div>
</template>

<style scoped lang="scss">
.componenttag {
    display: flex;
    align-items: center;
    height: 1.5rem;
    width: fit-content;
    border: solid 1px var(--color-complement-text);
    border-radius: 5px;
    padding: 0 4px;
    margin-right: 4px;

    span {
        font-family: var(--font-icon);
        color: var(--color-complement-text);
        margin-right: 4px;
    }

    p {
        color: var(--color-complement-text)
    }
}

.componenttag-fill {
    background-color: var(--color-complement-text);

    span,
    p {
        color: white
    }
}
</style>