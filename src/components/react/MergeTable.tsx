import { useMemo } from 'react';

interface CellConfig {
  content: React.ReactNode;
  colspan?: number;
  rowspan?: number;
  align?: 'left' | 'center' | 'right';
  header?: boolean;
  className?: string;
  style?: React.CSSProperties;
}

type CellInput = React.ReactNode | CellConfig;

interface MergeTableProps {
  data?: CellInput[][];
  columns?: (React.ReactNode | CellConfig)[];
  rows?: CellInput[][];
  headerRows?: number;
  bordered?: boolean;
  striped?: boolean;
  hoverable?: boolean;
  compact?: boolean;
  className?: string;
  style?: React.CSSProperties;
}

function isCellConfig(cell: CellInput): cell is CellConfig {
  return typeof cell === 'object' && cell !== null && 'content' in cell;
}

function normalizeCell(cell: CellInput): CellConfig {
  if (isCellConfig(cell)) {
    return cell;
  }
  return { content: cell };
}

export default function MergeTable({
  data,
  columns,
  rows,
  headerRows = 0,
  bordered = true,
  striped = false,
  hoverable = true,
  compact = false,
  className = '',
  style = {},
}: MergeTableProps) {
  const tableData = useMemo(() => {
    if (data) return data;
    const result: CellInput[][] = [];
    if (columns) {
      result.push(columns);
    }
    if (rows) {
      result.push(...rows);
    }
    return result;
  }, [data, columns, rows]);

  const { grid, colCount } = useMemo(() => {
    if (tableData.length === 0) return { grid: [], colCount: 0 };

    let maxCols = 0;
    tableData.forEach(row => {
      let cols = 0;
      row.forEach(cell => {
        const config = normalizeCell(cell);
        cols += config.colspan || 1;
      });
      maxCols = Math.max(maxCols, cols);
    });

    const occupied: boolean[][] = tableData.map(() => Array(maxCols).fill(false));
    const grid: { cell: CellConfig; rowSpan: number; colSpan: number }[][] = [];

    tableData.forEach((row, rowIndex) => {
      const gridRow: { cell: CellConfig; rowSpan: number; colSpan: number }[] = [];
      let colIndex = 0;

      row.forEach(cell => {
        while (colIndex < maxCols && occupied[rowIndex][colIndex]) {
          gridRow.push({ cell: { content: null }, rowSpan: 0, colSpan: 0 });
          colIndex++;
        }

        const config = normalizeCell(cell);
        const rowSpan = config.rowspan || 1;
        const colSpan = config.colspan || 1;

        for (let r = 0; r < rowSpan; r++) {
          for (let c = 0; c < colSpan; c++) {
            if (rowIndex + r < occupied.length && colIndex + c < maxCols) {
              occupied[rowIndex + r][colIndex + c] = true;
            }
          }
        }

        gridRow.push({ cell: config, rowSpan, colSpan });
        colIndex += colSpan;
      });

      while (colIndex < maxCols) {
        if (!occupied[rowIndex][colIndex]) {
          gridRow.push({ cell: { content: null }, rowSpan: 0, colSpan: 0 });
        }
        colIndex++;
      }

      grid.push(gridRow);
    });

    return { grid, colCount: maxCols };
  }, [tableData]);

  if (grid.length === 0) return null;

  const baseStyle: React.CSSProperties = {
    width: '100%',
    borderCollapse: 'collapse',
    fontSize: compact ? '0.875rem' : '1rem',
    ...style,
  };

  const getCellStyle = (cell: CellConfig, rowIndex: number): React.CSSProperties => {
    const isHeader = cell.header || rowIndex < headerRows;
    return {
      padding: compact ? '0.5rem 0.75rem' : '0.75rem 1rem',
      textAlign: cell.align || (isHeader ? 'center' : 'left'),
      fontWeight: isHeader ? 600 : 400,
      backgroundColor: cell.style?.backgroundColor,
      borderBottom: bordered ? '1px solid #e5e7eb' : undefined,
      borderRight: bordered ? '1px solid #e5e7eb' : undefined,
      borderLeft: bordered && rowIndex === 0 ? '1px solid #e5e7eb' : undefined,
      ...cell.style,
    };
  };

  return (
    <div 
      className={`merge-table-wrapper ${className}`}
      style={{ 
        overflowX: 'auto', 
        margin: '1rem 0',
        border: bordered ? '1px solid #e5e7eb' : undefined,
        borderRadius: '0.5rem',
      }}
    >
      <table style={baseStyle}>
        <tbody>
          {grid.map((row, rowIndex) => {
            const isStriped = striped && rowIndex >= headerRows && (rowIndex - headerRows) % 2 === 1;
            return (
              <tr 
                key={rowIndex}
                style={{
                  backgroundColor: isStriped ? '#f9fafb' : undefined,
                }}
                className={hoverable && rowIndex >= headerRows ? 'hover:bg-gray-50 transition-colors' : ''}
              >
                {row.map((cellData, colIndex) => {
                  if (cellData.rowSpan === 0 || cellData.colSpan === 0) {
                    return null;
                  }

                  const { cell, rowSpan, colSpan } = cellData;
                  const isHeader = cell.header || rowIndex < headerRows;
                  const CellTag = isHeader ? 'th' : 'td';

                  return (
                    <CellTag
                      key={colIndex}
                      rowSpan={rowSpan > 1 ? rowSpan : undefined}
                      colSpan={colSpan > 1 ? colSpan : undefined}
                      style={getCellStyle(cell, rowIndex)}
                      className={cell.className}
                    >
                      {cell.content}
                    </CellTag>
                  );
                })}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}

export function Cell({
  children,
  colspan,
  rowspan,
  align,
  header,
  className,
  style,
}: {
  children?: React.ReactNode;
  colspan?: number;
  rowspan?: number;
  align?: 'left' | 'center' | 'right';
  header?: boolean;
  className?: string;
  style?: React.CSSProperties;
}): CellConfig {
  return {
    content: children,
    colspan,
    rowspan,
    align,
    header,
    className,
    style,
  };
}
