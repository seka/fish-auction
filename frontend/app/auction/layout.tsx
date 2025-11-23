import Link from 'next/link';

export default function AuctionLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <div className="min-h-screen bg-orange-50">
            {/* Header */}
            <header className="bg-white shadow-sm border-b border-orange-100 sticky top-0 z-50">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
                    <div className="flex items-center">
                        <Link href="/" className="text-gray-500 hover:text-gray-700 mr-4 transition-colors font-medium">
                            &larr; 退場
                        </Link>
                        <h1 className="text-2xl font-bold text-orange-600 tracking-tight">セリ会場</h1>
                    </div>
                    <div className="text-sm text-gray-500 font-medium bg-orange-50 px-3 py-1 rounded-full">
                        Live Bidding System
                    </div>
                </div>
            </header>

            {/* Main Content */}
            <main>
                {children}
            </main>
        </div>
    );
}
