'use client';

import Link from 'next/link';

export default function AdminDashboard() {
    const menuItems = [
        {
            title: '漁師管理',
            description: '漁師の登録・一覧表示',
            href: '/admin/fishermen',
            icon: '👨‍🌾',
            color: 'indigo',
        },
        {
            title: '中買人管理',
            description: '中買人の登録・一覧表示',
            href: '/admin/buyers',
            icon: '👔',
            color: 'green',
        },
        {
            title: '出品管理',
            description: 'セリへの出品登録',
            href: '/admin/items',
            icon: '🐟',
            color: 'orange',
        },
        {
            title: '会場管理',
            description: 'セリ会場の登録・管理',
            href: '/admin/venues',
            icon: '🏢',
            color: 'blue',
        },
        {
            title: 'セリ管理',
            description: 'セリの作成・スケジュール管理',
            href: '/admin/auctions',
            icon: '📅',
            color: 'purple',
        },
        {
            title: '請求書発行',
            description: '落札後の請求書発行',
            href: '/invoice',
            icon: '💰',
            color: 'yellow',
        },
    ];

    const getColorClasses = (color: string) => {
        const colors: Record<string, { bg: string; hover: string; iconBg: string; iconText: string }> = {
            indigo: { bg: 'bg-indigo-50', hover: 'hover:bg-indigo-100', iconBg: 'bg-indigo-100', iconText: 'text-indigo-600' },
            green: { bg: 'bg-green-50', hover: 'hover:bg-green-100', iconBg: 'bg-green-100', iconText: 'text-green-600' },
            orange: { bg: 'bg-orange-50', hover: 'hover:bg-orange-100', iconBg: 'bg-orange-100', iconText: 'text-orange-600' },
            blue: { bg: 'bg-blue-50', hover: 'hover:bg-blue-100', iconBg: 'bg-blue-100', iconText: 'text-blue-600' },
            purple: { bg: 'bg-purple-50', hover: 'hover:bg-purple-100', iconBg: 'bg-purple-100', iconText: 'text-purple-600' },
            yellow: { bg: 'bg-yellow-50', hover: 'hover:bg-yellow-100', iconBg: 'bg-yellow-100', iconText: 'text-yellow-600' },
        };
        return colors[color] || colors.indigo;
    };

    return (
        <div className="max-w-7xl mx-auto p-6">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-800">管理ダッシュボード</h1>
                <p className="text-gray-600 mt-2">各管理メニューを選択してください</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {menuItems.map((item) => {
                    const colors = getColorClasses(item.color);
                    return (
                        <Link
                            key={item.href}
                            href={item.href}
                            className={`block p-6 rounded-xl border-2 border-transparent ${colors.bg} ${colors.hover} transition-all duration-200 hover:shadow-lg hover:scale-105`}
                        >
                            <div className="flex items-start space-x-4">
                                <div className={`p-3 rounded-lg ${colors.iconBg}`}>
                                    <span className="text-3xl">{item.icon}</span>
                                </div>
                                <div className="flex-1">
                                    <h3 className="text-lg font-bold text-gray-900 mb-1">{item.title}</h3>
                                    <p className="text-sm text-gray-600">{item.description}</p>
                                </div>
                            </div>
                        </Link>
                    );
                })}
            </div>

            <div className="mt-12 p-6 bg-blue-50 border border-blue-200 rounded-xl">
                <h2 className="text-lg font-bold text-blue-900 mb-2">📌 使い方</h2>
                <ol className="list-decimal list-inside space-y-1 text-sm text-blue-800">
                    <li>まず「会場管理」でセリを行う会場を登録します</li>
                    <li>「セリ管理」で開催日時を設定してセリを作成します</li>
                    <li>「漁師管理」「中買人管理」で参加者を登録します</li>
                    <li>「出品管理」で魚を登録してセリに出品します</li>
                    <li>セリ会場で入札が行われます</li>
                    <li>「請求書発行」で落札後の請求書を発行します</li>
                </ol>
            </div>
        </div>
    );
}
