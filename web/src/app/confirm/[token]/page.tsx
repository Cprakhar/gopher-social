'use client'

import { ServerURL } from "@/app/page"
import { Mail } from "lucide-react"
import { useParams, useRouter } from "next/navigation"

const ConfirmPage = () => {
  const { token } = useParams()
  const router = useRouter()
  const handleConfirm = async () => {
    const response = await fetch(`${ServerURL}/v1/users/activate/${token}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    if (response.ok) {
      router.push('/')
    } else {
      // TODO handle error
      alert('Error confirming email')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-100 to-indigo-200">
      <div className="bg-white shadow-xl rounded-xl p-8 max-w-md w-full flex flex-col items-center">
        <div className="mb-4 flex items-center justify-center">
          <Mail className="w-8 h-8 text-blue-600" />
        </div>
        <h1 className="text-2xl font-bold text-gray-800 mb-2">Confirm your email</h1>
        <p className="text-gray-600 mb-6 text-center">Click the button below to confirm your email address and activate your account.</p>
        <button
          onClick={handleConfirm}
          className="w-full py-3 px-6 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg shadow transition duration-150"
        >
          Confirm Email
        </button>
      </div>
    </div>
  )
}

export default ConfirmPage