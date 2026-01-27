import { useEffect, useState } from 'react';
import { AppShell, Button, Container, Group, Stack, Text, Title } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import Editors from './components/Editors';

const sampleXsl = `<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform" version="1.0"
  xmlns:fo="http://www.w3.org/1999/XSL/Format">
  <xsl:output method="xml" indent="yes"/>
  <xsl:template match="/">
    <fo:root>
      <fo:layout-master-set>
        <fo:simple-page-master master-name="A4" page-height="29.7cm" page-width="21cm" margin="2cm">
          <fo:region-body/>
        </fo:simple-page-master>
      </fo:layout-master-set>
      <fo:page-sequence master-reference="A4">
        <fo:flow flow-name="xsl-region-body">
          <fo:block font-size="18pt" font-weight="bold">Sample Document</fo:block>
          <fo:block margin-top="10pt">Hello from fopeditor!</fo:block>
          <fo:block margin-top="10pt">
            <xsl:text>Customer: </xsl:text>
            <xsl:value-of select="/invoice/customer"/>
          </fo:block>
        </fo:flow>
      </fo:page-sequence>
    </fo:root>
  </xsl:template>
</xsl:stylesheet>`;

const sampleXml = `<?xml version="1.0" encoding="UTF-8"?>
<invoice>
  <customer>Jane Doe</customer>
</invoice>`;

const STORAGE_KEYS = {
  xsl: 'fopeditor:xsl',
  xml: 'fopeditor:xml',
} as const;

const readStoredValue = (key: string, fallback: string) => {
  if (typeof window === 'undefined') {
    return fallback;
  }
  return window.localStorage.getItem(key) ?? fallback;
};

function App() {
  const [xsl, setXsl] = useState(() => readStoredValue(STORAGE_KEYS.xsl, sampleXsl));
  const [xml, setXml] = useState(() => readStoredValue(STORAGE_KEYS.xml, sampleXml));
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [layout, setLayout] = useState<'vertical' | 'horizontal'>('vertical');
  const [splitSizes, setSplitSizes] = useState({ vertical: [50, 50] as [number, number], horizontal: [50, 50] as [number, number] });

  useEffect(() => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(STORAGE_KEYS.xsl, xsl);
    }
  }, [xsl]);

  useEffect(() => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(STORAGE_KEYS.xml, xml);
    }
  }, [xml]);

  const handleRender = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch('/api/render', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ xsl, xml }),
      });
      if (!response.ok) {
        const message = (await response.text()).trim();
        throw new Error(message || 'Render failed');
      }
      const blob = await response.blob();
      const blobUrl = URL.createObjectURL(blob);
      window.open(blobUrl, '_blank');
      notifications.show({
        title: 'Render complete',
        message: 'PDF opened in a new tab',
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      setError(message);
      notifications.show({
        color: 'red',
        title: 'Render failed',
        message: message.length > 200 ? `${message.slice(0, 200)}…` : message,
      });
    } finally {
      setLoading(false);
    }
  };

  const handleExample = () => {
    setXsl(sampleXsl);
    setXml(sampleXml);
    setError(null);
  };

  return (
    <AppShell padding="md">
      <AppShell.Main>
        <Container fluid>
          <Stack gap="md">
            <Group justify="space-between" align="center">
              <div>
                <Title order={3} m={0}>
                  fopeditor
                </Title>
                <Text c="dimmed" size="sm">
                  Edit XSL-FO and XML, then render PDFs with Apache FOP.
                </Text>
              </div>
              <Group gap="xs">
                <Button onClick={handleRender} loading={loading}>
                  Render PDF
                </Button>
                <Button variant="default" onClick={handleExample} disabled={loading}>
                  Load example
                </Button>
                <Button
                  variant="light"
                  onClick={() => setLayout((prev) => (prev === 'vertical' ? 'horizontal' : 'vertical'))}
                  disabled={loading}
                >
                  {layout === 'vertical' ? 'Stack editors' : 'Split editors'}
                </Button>
              </Group>
            </Group>
            {error && (
              <Text
                c="red"
                size="sm"
                style={{
                  whiteSpace: 'pre-wrap',
                  fontFamily: 'var(--mantine-font-family-monospace)',
                  border: '1px solid var(--mantine-color-red-4)',
                  borderRadius: 'var(--mantine-radius-sm)',
                  padding: 'var(--mantine-spacing-xs)',
                  maxHeight: 240,
                  overflowY: 'auto',
                }}
              >
                {error}
              </Text>
            )}
            <Editors
              xsl={xsl}
              xml={xml}
              onXslChange={setXsl}
              onXmlChange={setXml}
              layout={layout}
              sizes={splitSizes[layout]}
              onResize={(next) =>
                setSplitSizes((prev) => ({
                  ...prev,
                  [layout]: next,
                }))
              }
              disabled={loading}
            />
          </Stack>
        </Container>
      </AppShell.Main>
    </AppShell>
  );
}

export default App;
