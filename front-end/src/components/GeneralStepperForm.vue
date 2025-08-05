<script setup lang="ts">
import {Check, Circle, Dot} from 'lucide-vue-next'
import {computed, ref} from 'vue'
import * as z from 'zod'
import {Button} from '@/components/ui/button'
import {
  Stepper,
  StepperDescription,
  StepperItem,
  StepperSeparator,
  StepperTitle,
  StepperTrigger
} from '@/components/ui/stepper'
import {useField, useForm} from "vee-validate";
import {toTypedSchema} from "@vee-validate/zod";
import PinInput from "@/components/PinInput.vue";
import {useAuthStore} from "@/stores/useAuth";
import {toast} from "vue-sonner";
import {useScreenSheetStore} from "@/stores/useScreenSheetStore";

const auth = useAuthStore()
const {setOpenDialog} = useScreenSheetStore()
const schema = z.object({
  email: z.string().min(1, "Email is required").email("Email is not valid"),
  password: z.string().min(1, "Password is required").min(6, "Password should be at least 6 characters"),
})

const {handleSubmit, meta, values} = useForm({
  validationSchema: toTypedSchema(schema),
})

const {
  value: email,
  errorMessage: emailError,
} = useField('email')
const {
  value: password,
  errorMessage: passwordError,
} = useField('password')

const stepIndex = ref(1)
const codeEntered = ref("")
const finalCode = computed(() => codeEntered.value)
const steps = [
  {
    step: 1,
    title: 'Your details',
    description: 'Provide your email and password',
  },
  {
    step: 2,
    title: 'Confirmation',
    description: 'Confirm your email',
  }
]

const currentStep = computed(() => stepIndex.value)

const hasFormData = computed(() => {
  return email.value && password.value
})

const isNextDisabled = computed(() => {
  if (currentStep.value >= steps.length) return true
  if (currentStep.value === 1) return !hasFormData.value || !meta.value.valid
  return false
})

const isPrevDisabled = computed(() => currentStep.value <= 1)

const handlePrev = () => {
  stepIndex.value--
}

const onSubmit = handleSubmit(async values => {
  console.log(values)
  try {
    const res = await auth.registration(values)

    if (res == 201) {
      toast.success("You've been registered. Please verify the email.")
      stepIndex.value++
    }
  } catch (error) {
    console.error(error)
    toast(auth.error)
  }
})

const handleCode = async () => {
  try {
    stepIndex.value++
    const res = await auth.verifyCode(finalCode.value, values.email!)

    if (res) {
      setOpenDialog(false, {name: 'RegisterDialog'})
    }
  } catch (error) {
    console.error(error)
  }
}

const handleResend = async () => {
  stepIndex.value = 2
  try {
    const res = await auth.resendCode(values.email!)

    if (res != 200) {
      stepIndex.value = 3
      console.error(auth.error)
    }
  } catch (error) {
    console.error(error)
    stepIndex.value = 3
  }
}
</script>

<template>
  <Stepper v-model="stepIndex" class="block w-full">
    <form
        @submit.prevent="onSubmit"
    >
      <div class="flex w-full flex-start gap-2 my-5 mb-9">
        <StepperItem
            v-for="step in steps"
            :key="step.step"
            v-slot="{ state }"
            class="relative flex w-full flex-col items-center justify-center"
            :step="step.step"
        >
          <StepperSeparator
              v-if="step.step !== steps[steps.length - 1].step"
              class="absolute left-[calc(50%+20px)] right-[calc(-50%+10px)] top-5 block h-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-primary"
          />

          <StepperTrigger as-child>
            <Button
                :variant="state === 'completed' || state === 'active' ? 'default' : 'outline'"
                size="icon"
                class="z-10 rounded-full shrink-0"
                :class="[state === 'active' && 'ring-2 ring-ring ring-offset-2 ring-offset-background']"
                :disabled="state !== 'completed'"
            >
              <Check v-if="state === 'completed'" class="size-5"/>
              <Circle v-if="state === 'active'"/>
              <Dot v-if="state === 'inactive'"/>
            </Button>
          </StepperTrigger>

          <div class="mt-5 flex flex-col items-center text-center">
            <StepperTitle
                :class="[state === 'active' && 'text-primary']"
                class="text-sm font-semibold transition lg:text-base"
            >
              {{ step.title }}
            </StepperTitle>
            <StepperDescription
                :class="[state === 'active' && 'text-primary']"
                class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm"
            >
              {{ step.description }}
            </StepperDescription>
          </div>
        </StepperItem>
      </div>

      <div class="flex flex-col gap-4 mt-4">
        <template v-if="stepIndex === 1">
          <div>
            <label for="email">Email</label>
            <input placeholder="example@email.com" id="email" type="email" v-model="email"
                   class="border rounded-sm px-2 py-1 w-full"/>
            <span class="text-red-500 text-sm">{{ emailError }}</span>
            <span class="text-red-500 text-sm">{{ auth.error }}</span>
          </div>

          <div>
            <label for="password">Password</label>
            <input placeholder="yourpassword" id="password" type="password" v-model="password"
                   class="rounded-sm border px-2 py-1 w-full"/>
            <span class="text-red-500 text-sm">{{ passwordError }}</span>
          </div>
        </template>
        <template v-if="stepIndex === 2">
          <div class="flex flex-col items-center justify-center gap-4 mb-5">
            <h1 class="font-semibold text-lg">Enter verification code</h1>
            <PinInput :set-code="(code) => codeEntered = code" class="w-fit"></PinInput>
            <span class="text-red-500 text-sm">{{ auth.error }}</span>
          </div>
        </template>
        <template v-if="stepIndex === 3">
          <div class="flex flex-col items-center justify-center gap-4 mb-5">
            <h1 v-if="auth.loading" class="font-semibold text-lg">Verifying...</h1>
            <div v-if="auth.error" class="flex flex-col items-center gap-2">
              <span class="text-red-500 text-sm">Failed to verify the email. Please try again.</span>
              <Button type="button" class="cursor-pointer" @click="() => handleResend()">Retry</Button>
            </div>
          </div>
        </template>
      </div>

      <div class="flex items-center justify-between mt-4">
        <Button class="cursor-pointer" :disabled="isPrevDisabled" variant="outline"
                size="sm" @click="handlePrev">
          Back
        </Button>
        <div class="flex items-center gap-3">
          <Button v-if="stepIndex === 1" :disabled="isNextDisabled" size="sm" type="submit"
                  class="cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]">
            Next
          </Button>
          <Button @click="handleCode" :disabled="finalCode.length < 5" v-if="stepIndex === 2" size="sm" type="button"
                  class="cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]">
            Verify
          </Button>
        </div>
      </div>
    </form>
  </Stepper>
</template>