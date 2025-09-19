<script setup>
import * as z from "zod";
import {useField, useForm} from "vee-validate";
import {toTypedSchema} from "@vee-validate/zod";
import {Button} from '@/components/ui/button'
import {useAuthStore} from "@/stores/useAuth.js";
import {toast} from "vue-sonner";
import {ref} from "vue";

const schema = z.object({
  email: z.string().min(1, "Email is required").email("Email is not valid").default(""),
})

const auth = useAuthStore()
const emailSent = ref(false)

const {handleSubmit, meta} = useForm({
  validationSchema: toTypedSchema(schema),
})

const {
  value: email,
  errorMessage: emailError,
} = useField('email')

const onSubmit = handleSubmit(async values => {
  try {
    const res = await auth.requestPasswordReset(values.email)

    if (res === 200) {
      toast.success("Password reset email sent successfully")
      emailSent.value = true
    } else if (res === 404) {
      toast.error("No account found with this email address")
    } else {
      toast.error("Failed to send reset email. Please try again.")
    }
  } catch (error) {
    console.error(error)
    toast.error("An unexpected error occurred")
  }
})

const emit = defineEmits(['back-to-login'])
</script>

<template>
  <div class="flex flex-col gap-8">
    <h1 class="font-semibold text-2xl text-center mt-4">FUMI</h1>
    
    <div v-if="!emailSent" class="flex flex-col gap-4">
      <h2 class="font-semibold text-xl text-center">Reset Password</h2>
      <p class="text-sm text-gray-600 text-center">
        Enter your email address and we'll send you a link to reset your password.
      </p>
      
      <form @submit.prevent="onSubmit" class="flex flex-col w-full items-center gap-4">
        <div class="flex flex-col w-full">
          <label for="reset-email">Email</label>
          <input 
            class="border rounded-sm px-2 py-1" 
            v-model="email" 
            id="reset-email" 
            placeholder="example@email.com"
            type="email"
          />
          <span class="text-red-500 text-sm">{{ emailError }}</span>
          <span class="text-red-500 text-sm">{{ auth.error }}</span>
        </div>
        
        <div class="flex justify-between items-center w-full">
          <Button type="button" class="cursor-pointer" variant="outline" @click="emit('back-to-login')">
            Back to Login
          </Button>
          <Button 
            type="submit" 
            class="cursor-pointer bg-[#c9a275] hover:bg-[#dbb384]"
            :disabled="!meta.valid || auth.loading"
          >
            {{ auth.loading ? 'Sending...' : 'Send Reset Link' }}
          </Button>
        </div>
      </form>
    </div>

    <div v-else class="flex flex-col gap-4 text-center">
      <h2 class="font-semibold text-xl">Email Sent!</h2>
      <p class="text-sm text-gray-600">
        We've sent a password reset link to <strong>{{ email }}</strong>. 
        Please check your email and follow the instructions to reset your password.
      </p>
      <p class="text-xs text-gray-500">
        Didn't receive the email? Check your spam folder or try again.
      </p>
      
      <div class="flex justify-center mt-4">
        <Button 
          type="button" 
          class="cursor-pointer" 
          variant="outline" 
          @click="emit('back-to-login')"
        >
          Back to Login
        </Button>
      </div>
    </div>
  </div>
</template>