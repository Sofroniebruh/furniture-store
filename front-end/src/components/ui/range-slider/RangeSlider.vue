<template>
  <div
      ref="sliderRef"
      class="relative flex w-full touch-none select-none mb-6 items-center"
      :class="className"
      @mousedown="handleMouseDown"
      @touchstart="handleTouchStart"
  >
    <!-- Track -->
    <div class="relative h-1 w-full grow overflow-hidden rounded-full bg-gray-100">
      <!-- Range -->
      <div
          class="absolute h-full bg-[#c9a275]"
          :style="{
          left: `${((localValues[0] - min) / (max - min)) * 100}%`,
          width: `${((localValues[1] - localValues[0]) / (max - min)) * 100}%`
        }"
      />
    </div>

    <!-- Thumbs and Labels -->
    <div
        v-for="(value, index) in localValues"
        :key="index"
        class="absolute"
        :style="{
        left: `calc(${((value - min) / (max - min)) * 100}% - 8px)`
      }"
    >
      <!-- Label -->
      <div
          class="absolute text-center"
          :style="{
          left: '8px',
          top: '10px',
          transform: 'translateX(-50%)'
        }"
      >
        <span class="text-sm">{{ formatLabel ? formatLabel(value) : value }}</span>
      </div>

      <!-- Thumb -->
      <div
          class="block h-4 w-4 rounded-full border border-primary/50 bg-white shadow transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 cursor-pointer"
          :data-thumb-index="index"
          @mousedown="(e) => handleThumbMouseDown(e, index)"
          @touchstart="(e) => handleThumbTouchStart(e, index)"
          @keydown="(e) => handleKeyDown(e, index)"
          tabindex="0"
          role="slider"
          :aria-valuemin="min"
          :aria-valuemax="max"
          :aria-valuenow="value"
          :aria-label="`Slider thumb ${index + 1}`"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import {computed, onMounted, onUnmounted, ref, watch} from 'vue'

interface SliderProps {
  className?: string
  min: number
  max: number
  step: number
  formatLabel?: (value: number) => string
  modelValue?: number[] | readonly number[]
  onValueChange?: (values: number[]) => void
}

const props = withDefaults(defineProps<SliderProps>(), {
  className: '',
  modelValue: undefined,
  onValueChange: undefined
})

const emit = defineEmits<{
  'update:modelValue': [value: number[]]
  'valueChange': [value: number[]]
}>()

const sliderRef = ref<HTMLDivElement>()
const activeThumbIndex = ref<number | null>(null)
const isDragging = ref(false)

const initialValue = computed(() => {
  if (Array.isArray(props.modelValue)) {
    return props.modelValue
  }
  return [props.min, props.max]
})

const localValues = ref<number[]>(initialValue.value)

// Watch for external value changes
watch(
    () => props.modelValue,
    (newValue) => {
      if (Array.isArray(newValue)) {
        localValues.value = [...newValue]
      } else {
        localValues.value = [props.min, props.max]
      }
    },
    {immediate: true}
)

const handleValueChange = (newValues: number[]) => {
  localValues.value = [...newValues]
  emit('update:modelValue', newValues)
  if (props.onValueChange) {
    props.onValueChange(newValues)
  }
  emit('valueChange', newValues)
}

const snapToStep = (value: number): number => {
  const steps = Math.round((value - props.min) / props.step)
  return Math.max(props.min, Math.min(props.max, props.min + steps * props.step))
}

const getValueFromPosition = (clientX: number): number => {
  if (!sliderRef.value) return props.min

  const rect = sliderRef.value.getBoundingClientRect()
  const percentage = Math.max(0, Math.min(1, (clientX - rect.left) / rect.width))
  const value = props.min + percentage * (props.max - props.min)
  return snapToStep(value)
}

const updateValueAtIndex = (index: number, newValue: number) => {
  const newValues = [...localValues.value]
  newValues[index] = newValue

  // Ensure values don't cross over
  if (index === 0 && newValue > newValues[1]) {
    newValues[0] = newValues[1]
  } else if (index === 1 && newValue < newValues[0]) {
    newValues[1] = newValues[0]
  } else {
    newValues[index] = newValue
  }

  handleValueChange(newValues)
}

const handleMouseMove = (e: MouseEvent) => {
  if (!isDragging.value || activeThumbIndex.value === null) return

  const newValue = getValueFromPosition(e.clientX)
  updateValueAtIndex(activeThumbIndex.value, newValue)
}

const handleMouseUp = () => {
  isDragging.value = false
  activeThumbIndex.value = null
}

const handleTouchMove = (e: TouchEvent) => {
  if (!isDragging.value || activeThumbIndex.value === null) return

  const touch = e.touches[0]
  const newValue = getValueFromPosition(touch.clientX)
  updateValueAtIndex(activeThumbIndex.value, newValue)
}

const handleTouchEnd = () => {
  isDragging.value = false
  activeThumbIndex.value = null
}

const handleThumbMouseDown = (e: MouseEvent, index: number) => {
  e.preventDefault()
  e.stopPropagation()
  isDragging.value = true
  activeThumbIndex.value = index
}

const handleThumbTouchStart = (e: TouchEvent, index: number) => {
  e.preventDefault()
  e.stopPropagation()
  isDragging.value = true
  activeThumbIndex.value = index
}

const handleMouseDown = (e: MouseEvent) => {
  if (isDragging.value) return

  const newValue = getValueFromPosition(e.clientX)
  const distances = localValues.value.map(val => Math.abs(val - newValue))
  const closestIndex = distances.indexOf(Math.min(...distances))

  updateValueAtIndex(closestIndex, newValue)
}

const handleTouchStart = (e: TouchEvent) => {
  if (isDragging.value) return

  const touch = e.touches[0]
  const newValue = getValueFromPosition(touch.clientX)
  const distances = localValues.value.map(val => Math.abs(val - newValue))
  const closestIndex = distances.indexOf(Math.min(...distances))

  updateValueAtIndex(closestIndex, newValue)
}

const handleKeyDown = (e: KeyboardEvent, index: number) => {
  let newValue = localValues.value[index]

  switch (e.key) {
    case 'ArrowLeft':
    case 'ArrowDown':
      e.preventDefault()
      newValue = Math.max(props.min, newValue - props.step)
      break
    case 'ArrowRight':
    case 'ArrowUp':
      e.preventDefault()
      newValue = Math.min(props.max, newValue + props.step)
      break
    case 'Home':
      e.preventDefault()
      newValue = props.min
      break
    case 'End':
      e.preventDefault()
      newValue = props.max
      break
    case 'PageUp':
      e.preventDefault()
      newValue = Math.min(props.max, newValue + props.step * 10)
      break
    case 'PageDown':
      e.preventDefault()
      newValue = Math.max(props.min, newValue - props.step * 10)
      break
  }

  updateValueAtIndex(index, newValue)
}

onMounted(() => {
  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
  document.addEventListener('touchmove', handleTouchMove)
  document.addEventListener('touchend', handleTouchEnd)
})

onUnmounted(() => {
  document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseup', handleMouseUp)
  document.removeEventListener('touchmove', handleTouchMove)
  document.removeEventListener('touchend', handleTouchEnd)
})
</script>