import Link from 'next/link';

export default function AdminLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <div className="flex min-h-screen bg-gray-100">
            {/* Sidebar */}
            <aside className="w-64 bg-indigo-900 text-white flex-shrink-0 shadow-xl">
                <div className="p-6 bg-indigo-950">
                    <h2 className="text-xl font-bold tracking-wider">管理画面</h2>
                    <p className="text-xs text-indigo-300 mt-1">Fish Auction Admin</p>
                </div>
                <nav className="mt-6 px-2 space-y-1">
                    <Link href="/" className="block py-3 px-4 rounded hover:bg-indigo-800 transition-colors text-sm font-medium">
                        &larr; トップに戻る
                    </Link>
                    <div className="border-t border-indigo-800 my-4 mx-2"></div>
                    <Link href="/admin" className="block py-3 px-4 rounded hover:bg-indigo-800 transition-colors text-sm font-medium">
                        マスタ・出品登録
                    </Link>
                    <Link href="/invoice" className="block py-3 px-4 rounded hover:bg-indigo-800 transition-colors text-sm font-medium">
                        請求書発行
                    </Link>
                </nav>
            </aside>

            {/* Main Content */}
            <main className="flex-1 p-8 overflow-y-auto">
                {children}
            </main>
        </div>
    );
}
