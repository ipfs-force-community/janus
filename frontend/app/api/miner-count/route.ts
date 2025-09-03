import { NextResponse } from 'next/server';

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const range = searchParams.get('range') || '30d';

  try {
    const res = await fetch(`http://127.0.0.1:10086/miners?interval=${range}`, {
      method: 'GET',
    });

    if (!res.ok) {
      return NextResponse.json(
        { error: `Backend error: ${res.status} ${res.statusText}` },
        { status: res.status },
      );
    }

    const data = await res.json();
    return NextResponse.json(data);
  } catch (err: any) {
    return NextResponse.json(
      { error: `Failed to fetch backend: ${err.message}` },
      { status: 500 },
    );
  }
}
