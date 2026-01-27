import { Paper, Text } from '@mantine/core';
import Editor from '@monaco-editor/react';
import Split from 'react-split';
import './Editors.css';

interface EditorsProps {
  xsl: string;
  xml: string;
  onXslChange: (value: string) => void;
  onXmlChange: (value: string) => void;
  layout: 'vertical' | 'horizontal';
  sizes: [number, number];
  onResize: (sizes: [number, number]) => void;
  disabled?: boolean;
}

const editorOptions = {
  fontSize: 14,
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
};

function Editors({ xsl, xml, onXslChange, onXmlChange, layout, sizes, onResize, disabled }: EditorsProps) {
  const direction = layout === 'vertical' ? 'horizontal' : 'vertical';
  const splitClass = layout === 'vertical' ? 'editors-split' : 'editors-split editors-split-horizontal';

  const panel = (label: string, value: string, onChange: (value: string) => void, path: string) => (
    <div className="editors-panel">
      <Text fw={500}>{label}</Text>
      <div className="editors-editor">
        <Paper shadow="sm" radius="md" withBorder className="editors-paper">
          <Editor
            height="100%"
            width="100%"
            language="xml"
            theme="vs-dark"
            value={value}
            onChange={(val) => onChange(val ?? '')}
            options={{ ...editorOptions, readOnly: Boolean(disabled), automaticLayout: true }}
            loading="Loading editor..."
            path={path}
            saveViewState
            keepCurrentModel
          />
        </Paper>
      </div>
    </div>
  );

  return (
    <div className="editors-container">
      <Split
        key={layout}
        className={splitClass}
        direction={direction}
        gutterSize={12}
        minSize={layout === 'vertical' ? 320 : 200}
        sizes={sizes}
        onDragEnd={(next) => onResize([next[0], next[1]])}
        snapOffset={0}
        style={{ width: '100%', height: '100%' }}
      >
        {panel('XSL / XSL-FO', xsl, onXslChange, 'xsl.xml')}
        {panel('XML data', xml, onXmlChange, 'xml.xml')}
      </Split>
    </div>
  );
}

export default Editors;
