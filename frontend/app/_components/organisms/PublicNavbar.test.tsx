import { render, screen, fireEvent } from "@testing-library/react";
import { PublicNavbar } from "./PublicNavbar";
import { vi, describe, it, expect, beforeEach } from "vitest";
import { usePathname, useRouter } from "next/navigation";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import * as buyerAuth from '@/src/api/buyer_auth';

// Mocks
vi.mock("next/navigation", () => ({
    usePathname: vi.fn(),
    useRouter: vi.fn(),
}));

vi.mock("next-intl", () => ({
    useTranslations: () => (key: string) => key,
}));

vi.mock("@tanstack/react-query", () => ({
    useQuery: vi.fn(),
    useQueryClient: vi.fn(),
}));

vi.mock("@/src/api/buyer_auth", () => ({
    getCurrentBuyer: vi.fn(),
    logoutBuyer: vi.fn(),
}));

describe("PublicNavbar", () => {
    const mockPush = vi.fn();
    const mockSetQueryData = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        (useRouter as any).mockReturnValue({ push: mockPush });
        (useQueryClient as any).mockReturnValue({ setQueryData: mockSetQueryData });
        (usePathname as any).mockReturnValue("/");
    });

    it("renders correctly when not logged in", () => {
        (useQuery as any).mockReturnValue({ data: null });

        render(<PublicNavbar />);

        expect(screen.getByText("Common.app_name")).toBeInTheDocument();
        expect(screen.getByText("Navbar.active_auctions")).toBeInTheDocument();
        expect(screen.getByText("Navbar.login")).toBeInTheDocument();
        expect(screen.queryByText("Navbar.logout")).not.toBeInTheDocument();
        expect(screen.queryByText("Navbar.mypage")).not.toBeInTheDocument();
    });

    it("renders correctly when logged in", () => {
        (useQuery as any).mockReturnValue({ data: { id: 1, name: "Buyer" } });

        render(<PublicNavbar />);

        expect(screen.getByText("Navbar.logout")).toBeInTheDocument();
        expect(screen.getByText("Navbar.mypage")).toBeInTheDocument();
        expect(screen.queryByText("Navbar.login")).not.toBeInTheDocument();
    });

    it("handles logout correctly", async () => {
        (useQuery as any).mockReturnValue({ data: { id: 1, name: "Buyer" } });
        const mockLogoutBuyer = vi.spyOn(buyerAuth, 'logoutBuyer').mockResolvedValue(undefined);

        render(<PublicNavbar />);

        const logoutButton = screen.getByText("Navbar.logout");
        fireEvent.click(logoutButton);

        expect(mockLogoutBuyer).toHaveBeenCalled();
        // Wait for async actions if necessary, but here spy serves.
        // In real execution, we might need waitFor if state update is involved,
        // but here handleLogout calls router.push immediately after await.
        // Since logoutBuyer is mocked to resolve immediately:
    });

    it("does not render on admin pages", () => {
        (usePathname as any).mockReturnValue("/admin/dashboard");
        const { container } = render(<PublicNavbar />);
        expect(container).toBeEmptyDOMElement();
    });
});
