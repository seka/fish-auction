'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { buyerSignupSchema, BuyerSignupFormData } from '@/src/models/schemas/buyer_auth';
import { signupBuyer } from '@/src/api/buyer_auth';
import Link from 'next/link';

export default function SignupPage() {
    const [error, setError] = useState('');
    const router = useRouter();
    const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<BuyerSignupFormData>({
        resolver: zodResolver(buyerSignupSchema),
    });

    const onSubmit = async (data: BuyerSignupFormData) => {
        setError('');
        try {
            await signupBuyer(data);
            router.push('/login/buyer');
        } catch (e: any) {
            if (e.response && e.response.status === 409) {
                setError('登録に失敗しました。名前が既に使用されている可能性があります。');
            } else if (e.response && e.response.status >= 500) {
                setError('この操作の実行中にエラーが発生しました。運営にお問い合わせください');
            } else {
                setError('登録に失敗しました。入力内容をご確認ください。');
            }
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
            <div className="max-w-md w-full space-y-8">
                <div>
                    <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
                        中買人登録
                    </h2>
                    <p className="mt-2 text-center text-sm text-gray-600">
                        セリに参加するにはアカウントを作成してください
                    </p>
                </div>
                <form className="mt-8 space-y-6" onSubmit={handleSubmit(onSubmit)}>
                    <div className="rounded-md shadow-sm -space-y-px">
                        <div>
                            <label htmlFor="name" className="sr-only">名前</label>
                            <input
                                id="name"
                                type="text"
                                {...register('name')}
                                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                placeholder="名前"
                            />
                            {errors.name && <p className="text-red-500 text-xs mt-1">{errors.name.message}</p>}
                        </div>
                        <div>
                            <label htmlFor="email" className="sr-only">メールアドレス</label>
                            <input
                                id="email"
                                type="email"
                                {...register('email')}
                                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                placeholder="メールアドレス"
                            />
                            {errors.email && <p className="text-red-500 text-xs mt-1">{errors.email.message}</p>}
                        </div>
                        <div>
                            <label htmlFor="organization" className="sr-only">所属組織</label>
                            <input
                                id="organization"
                                type="text"
                                {...register('organization')}
                                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                placeholder="所属組織"
                            />
                            {errors.organization && <p className="text-red-500 text-xs mt-1">{errors.organization.message}</p>}
                        </div>
                        <div>
                            <label htmlFor="contact_info" className="sr-only">連絡先</label>
                            <input
                                id="contact_info"
                                type="text"
                                {...register('contact_info')}
                                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                placeholder="連絡先"
                            />
                            {errors.contact_info && <p className="text-red-500 text-xs mt-1">{errors.contact_info.message}</p>}
                        </div>
                        <div>
                            <label htmlFor="password" className="sr-only">パスワード</label>
                            <input
                                id="password"
                                type="password"
                                {...register('password')}
                                className="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                                placeholder="パスワード"
                            />
                            {errors.password && <p className="text-red-500 text-xs mt-1">{errors.password.message}</p>}
                        </div>
                    </div>

                    {error && <div className="text-red-500 text-sm text-center">{error}</div>}

                    <div>
                        <button
                            type="submit"
                            disabled={isSubmitting}
                            className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                        >
                            登録
                        </button>
                    </div>
                    <div className="text-center">
                        <Link href="/login/buyer" className="text-sm text-indigo-600 hover:text-indigo-500">
                            すでにアカウントをお持ちの方はこちら
                        </Link>
                    </div>
                </form>
            </div>
        </div>
    );
}
