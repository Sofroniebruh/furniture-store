<script setup lang="ts">
import {ref, watch} from 'vue'
import {PinInput, PinInputGroup, PinInputSeparator, PinInputSlot,} from '@/components/ui/pin-input'

const value = ref<string[]>([])
const timeout = ref<number | null>(null)

const props = defineProps<{
  setCode: (code: string) => void
}>()

watch(value, (newVal) => {
  if (timeout.value) clearTimeout(timeout.value)

  timeout.value = setTimeout(() => {
    props.setCode(newVal.join(""))
  }, 500)
})
</script>

<template>
  <div class="w-fit">
    <PinInput
        id="pin-input"
        v-model="value"
        placeholder="○"
    >
      <PinInputGroup class="gap-1">
        <template v-for="(id, index) in 5" :key="id">
          <PinInputSlot
              class="rounded-md border"
              :index="index"
          />
          <template v-if="index !== 4">
            <PinInputSeparator/>
          </template>
        </template>
      </PinInputGroup>
    </PinInput>
  </div>
</template>