<!-- Cleaned -->

<!-- This component has two modes "outline" and "fill" which is controlled via the prop "mode" -->

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
    min-height: 1.5rem;
    max-height: 1.5rem;
    width: fit-content;
    max-width: 88px;
    display: flex;
    align-items: center;
    margin-right: 4px;
    padding: 0 4px;
    border: solid 1px var(--color-complement-text);
    border-radius: 5px;
    overflow-x: hidden;
    pointer-events: none;
    user-select: none;

    span {
        margin-right: 4px;
        color: var(--color-complement-text);
        font-family: var(--font-icon);
    }

    p {
        max-height: 1rem;
        color: var(--color-complement-text);
    }

    &-fill {
        background-color: var(--color-complement-text);

        span,
        p {
            color: white
        }
    }
}
</style>