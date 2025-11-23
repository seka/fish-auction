import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-8 font-sans bg-gradient-to-br from-blue-50 to-indigo-50">
      <div className="text-center mb-16">
        <h1 className="text-5xl font-bold text-indigo-900 mb-6 tracking-tight">
          漁港のせりシステム
        </h1>
        <p className="text-xl text-gray-600">
          役割を選択してログインしてください
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 w-full max-w-5xl">
        {/* Admin Portal */}
        <Link href="/admin" className="group block p-10 bg-white border border-gray-100 rounded-3xl shadow-xl hover:shadow-2xl hover:-translate-y-1 transition-all duration-300 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-50 rounded-bl-full -mr-8 -mt-8 transition-transform group-hover:scale-110"></div>
          <div className="flex flex-col items-center text-center space-y-6 relative z-10">
            <div className="p-5 bg-indigo-100 rounded-2xl group-hover:bg-indigo-600 transition-colors duration-300 shadow-sm">
              <svg className="w-14 h-14 text-indigo-600 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div>
              <h2 className="text-3xl font-bold text-gray-900 mb-2">管理画面</h2>
              <p className="text-gray-500 leading-relaxed">
                漁師・中買人の登録、<br />出品登録、請求書の管理を行います。
              </p>
            </div>
            <span className="inline-flex items-center text-indigo-600 font-bold group-hover:translate-x-1 transition-transform">
              管理画面へ移動 <span className="ml-2">&rarr;</span>
            </span>
          </div>
        </Link>

        {/* Auction Floor */}
        <Link href="/auction" className="group block p-10 bg-white border border-gray-100 rounded-3xl shadow-xl hover:shadow-2xl hover:-translate-y-1 transition-all duration-300 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-32 h-32 bg-orange-50 rounded-bl-full -mr-8 -mt-8 transition-transform group-hover:scale-110"></div>
          <div className="flex flex-col items-center text-center space-y-6 relative z-10">
            <div className="p-5 bg-orange-100 rounded-2xl group-hover:bg-orange-600 transition-colors duration-300 shadow-sm">
              <svg className="w-14 h-14 text-orange-600 group-hover:text-white transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
              </svg>
            </div>
            <div>
              <h2 className="text-3xl font-bold text-gray-900 mb-2">セリ会場</h2>
              <p className="text-gray-500 leading-relaxed">
                出品されている魚の確認、<br />リアルタイムでの入札を行います。
              </p>
            </div>
            <span className="inline-flex items-center text-orange-600 font-bold group-hover:translate-x-1 transition-transform">
              会場へ入場 <span className="ml-2">&rarr;</span>
            </span>
          </div>
        </Link>
      </div>
    </div>
  );
}
