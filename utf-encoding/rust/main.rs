const MAX_ONE_BYTE: u32 = 0x80; // 128
const MAX_TWO_BYTE: u32 = 0x800; // 2048
const MAX_THREE_BYTE: u32 = 0x10000; // 65536

const MASK: u32 = 0x3F; // 63 // 00111111
const CONTINUATION_MASK: u32 = 0x80; // 128 // 10000000
const TWO_BYTE_MASK: u32 = 0xC0; // 192 // 11000000
const THREE_BYTE_MASK: u32 = 0xE0; // 224 // 11100000
const FOUR_BYTE_MASK: u32 = 0xF0; // 240 // 11110000

const MAX_UTF16: u32 = 0x10000; // 65536
const MIN_HIGH_SURROGATE: u32 = 0xD800; // 55.296
const MAX_HIGH_SURROGATE: u32 = 0xDBFF; // 56.319

const MIN_LOW_SURROGATE: u32 = 0xDC00; // 56.320
const MAX_LOW_SURROGATE: u32 = 0xDFFF; // 57.343

const SURROGATE_MASK: u32 = 0x3FF; // 1023 // 1111111111


fn encode_utf8(c: char, out: &mut Vec<u8>) -> Result<(), String> {
    let c_decimal: u32 = c as u32;

    if c_decimal < MAX_ONE_BYTE {
        out.push(c_decimal as u8);
        return Ok(());
    } 

    if c_decimal < MAX_TWO_BYTE {
        let b_one: u8 = ((c_decimal >> 6) | TWO_BYTE_MASK) as u8;
        let b_two: u8 = ((c_decimal & MASK) | CONTINUATION_MASK) as u8;
        out.push(b_one);
        out.push(b_two);
        return Ok(());
    }

    if c_decimal < MAX_THREE_BYTE {
        let b_one: u8 = ((c_decimal >> 12) | THREE_BYTE_MASK) as u8;
        let b_two: u8 = (((c_decimal >> 6) & MASK) | CONTINUATION_MASK) as u8;
        let b_three: u8 = ((c_decimal & MASK) | CONTINUATION_MASK) as u8;
        out.push(b_one);
        out.push(b_two);
        out.push(b_three);
        return Ok(());
    }

    let b_one: u8 = ((c_decimal >> 18) | FOUR_BYTE_MASK) as u8;
    let b_two: u8 = (((c_decimal >> 12) & MASK) | CONTINUATION_MASK) as u8;
    let b_three: u8 = (((c_decimal >> 6) & MASK) | CONTINUATION_MASK) as u8;
    let b_four: u8 = ((c_decimal & MASK) | CONTINUATION_MASK) as u8;
    out.push(b_one);
    out.push(b_two);
    out.push(b_three);
    out.push(b_four);

    Ok(())
} 

fn encode_utf16(c: char, out: &mut Vec<u16>) -> Result<(), String> {
    let c_decimal = c as u32;

    if c_decimal < MAX_UTF16 {
        out.push(c_decimal as u16);
        return Ok(());
    }

    let c_remainder = c_decimal - MAX_UTF16;
    let high = ((c_remainder >> 10) + MIN_HIGH_SURROGATE) as u16;
    let low = ((c_remainder & SURROGATE_MASK) + MIN_LOW_SURROGATE) as u16;

    out.push(high);
    out.push(low);

    return Ok(())

}

fn main() {
    let sigma = 'Ʃ';
    let star = '✪';
    let rose_emoji = '🌹';
    let poop = '💩';

    let mut res: Vec<u16> = Vec::new();

    if let Err(e) = encode_utf16(poop, &mut res) {
        println!("{}", e);
        std::process::exit(1);
    }

    println!("{:?}", res);
}
