<script setup lang="ts">
import * as z from "zod";
import {useField, useForm} from "vee-validate";
import {toTypedSchema} from "@vee-validate/zod";
import {Button} from '@/components/ui/button'
import {useAuthStore} from "@/stores/useAuth.js";
import {toast} from "vue-sonner";
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref} from "vue";

const schema = z.object({
  password: z.string().min(1, "Password is required").min(6, "Password should be at least 6 characters").default(""),
  confirmPassword: z.string().min(1, "Please confirm your password").default(""),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
})

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const token = ref("")
const resetSuccessful = ref(false)

const {handleSubmit, meta} = useForm({
  validationSchema: toTypedSchema(schema),
})

const {
  value: password,
  errorMessage: passwordError,
} = useField('password')

const {
  value: confirmPassword,
  errorMessage: confirmPasswordError,
} = useField('confirmPassword')

onMounted(() => {
  token.value = route.query.token as string || ""
  if (!token.value) {
    toast.error("Invalid reset link")
    router.push('/')
  }
})

const onSubmit = handleSubmit(async values => {
  try {
    const res = await auth.resetPassword(token.value, values.password)

    if (res === 200) {
      toast.success("Password reset successfully!")
      resetSuccessful.value = true
      setTimeout(() => {
        router.push('/')
      }, 2000)
    } else if (res === 400) {
      toast.error("Invalid or expired reset token")
    } else {
      toast.error("Failed to reset password. Please try again.")
    }
  } catch (error) {
    console.error(error)
    toast.error("An unexpected error occurred")
  }
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="flex flex-col gap-8">
        <h1 class="font-semibold text-2xl text-center mt-4">FUMI</h1>
        
        <div v-if="!resetSuccessful" class="flex flex-col gap-4">
          <h2 class="font-semibold text-xl text-center">Set New Password</h2>
          <p class="text-sm text-gray-600 text-center">
            Enter your new password below.
          </p>
          
          <form @submit.prevent="onSubmit" class="flex flex-col w-full items-center gap-4">
            <div class="flex flex-col w-full">
              <label for="new-password">New Password</label>
              <input 
                class="border rounded-sm px-2 py-1" 
                v-model="password" 
                id="new-password" 
                placeholder="Enter new password"
                type="password"
              />
              <span class="text-red-500 text-sm">{{ passwordError }}</span>
            </div>
            
            <div class="flex flex-col w-full">
              <label for="confirm-password">Confirm Password</label>
              <input 
                class="border rounded-sm px-2 py-1" 
                v-model="confirmPassword" 
                id="confirm-password" 
                placeholder="Confirm new password"
                type="password"
              />
              <span class="text-red-500 text-sm">{{ confirmPasswordError }}</span>
              <span class="text-red-500 text-sm">{{ auth.error }}</span>
            </div>
            
            <div class="flex justify-center w-full">
              <Button 
                type="submit" 
                class="cursor-pointer bg-[#c9a275] hover:bg-[#dbb384] w-full"
                :disabled="!meta.valid || auth.loading"
              >
                {{ auth.loading ? 'Resetting...' : 'Reset Password' }}
              </Button>
            </div>
          </form>
        </div>

        <div v-else class="flex flex-col gap-4 text-center">
          <h2 class="font-semibold text-xl text-green-600">Password Reset Successful!</h2>
          <p class="text-sm text-gray-600">
            Your password has been successfully reset. You will be redirected to the login page shortly.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>