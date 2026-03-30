'use client';

import { Tr, Td } from '@molecules';
import { css } from 'styled-system/css';
import { InvoiceItem } from '@entities';

interface InvoiceListCellProps {
  invoice: InvoiceItem;
  onSelect: (invoice: InvoiceItem) => void;
  // t is actually unused for now, so I'll remove it from props but keep it in mind if needed
}

export const InvoiceListCell = ({ invoice, onSelect }: InvoiceListCellProps) => {
  return (
    <Tr onClick={() => onSelect(invoice)} _hover={{ bg: 'gray.50' }} transition="colors">
      <Td className={css({ fontSize: 'sm', color: 'gray.500', fontFamily: 'mono' })}>
        {invoice.buyerId}
      </Td>
      <Td className={css({ fontSize: 'sm', fontWeight: 'bold', color: 'gray.900' })}>
        {invoice.buyerName}
      </Td>
      <Td
        className={css({
          textAlign: 'right',
          fontWeight: 'bold',
          color: 'indigo.700',
          fontSize: 'lg',
        })}
      >
        ¥{invoice.totalAmount.toLocaleString()}
      </Td>
    </Tr>
  );
};
