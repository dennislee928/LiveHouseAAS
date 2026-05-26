"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { getToken, api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface Venue {
  id: string;
  name: string;
  city: string;
  capacity: number;
  status: string;
}

interface User {
  role: string;
}

export default function VenuesPage() {
  const [venues, setVenues] = useState<Venue[]>([]);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) return;

    api.get<User>("/api/v1/me", token).then(setUser);

    const endpoint = user?.role === "venue" ? "/api/v1/venues" : "/api/v1/venues/all";
    api.get<Venue[]>(endpoint, token)
      .then(setVenues)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [user?.role]);

  if (loading) return <p className="text-gray-500">載入中...</p>;

  const isVenueOwner = user?.role === "venue";

  return (
    <div>
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">場館</h2>
          <p className="mt-1 text-gray-600">
            {isVenueOwner ? "管理您的場館資訊" : "瀏覽所有可用場館"}
          </p>
        </div>
        {isVenueOwner && (
          <Link href="/venues/new">
            <Button>新增場館</Button>
          </Link>
        )}
      </div>

      <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {venues.map((venue) => (
          <Link key={venue.id} href={`/venues/${venue.id}`}>
            <Card className="cursor-pointer transition-shadow hover:shadow-md">
              <CardHeader>
                <CardTitle className="text-lg">{venue.name}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-gray-500">{venue.city}</p>
                <p className="text-sm text-gray-500">容納人數: {venue.capacity}</p>
                <span className={`mt-2 inline-block rounded-full px-2 py-0.5 text-xs font-medium ${
                  venue.status === "active"
                    ? "bg-green-50 text-green-700"
                    : "bg-gray-50 text-gray-600"
                }`}>
                  {venue.status === "active" ? "營業中" : venue.status}
                </span>
              </CardContent>
            </Card>
          </Link>
        ))}
        {venues.length === 0 && (
          <p className="col-span-full text-center text-gray-400">尚無場館資料</p>
        )}
      </div>
    </div>
  );
}
