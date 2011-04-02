package main

import "fmt"
import "os"
import "io"
import "encoding/binary"

type vm struct {
    i int
}

const (
    ACC_PUBLIC      = 0x0001
    ACC_FINAL       = 0x0010
    ACC_SUPER       = 0x0020
    ACC_INTERFACE   = 0x0200
    ACC_ABSTRACT    = 0x0400
)

const (
    CONST_Class                 = 7
    CONST_FieldRef              = 9
    CONST_MethodRef             = 10
    CONST_InterfaceMethodRef    = 11
    CONST_String                = 8
    CONST_Integer               = 3
    CONST_Float                 = 4
    CONST_Long                  = 5
    CONST_Double                = 6
    CONST_NameAndType           = 12
    CONST_Utf8                  = 1
)

type cp_info struct {
    tag     uint8
    info    []uint8
}

type cp_class_info struct {
    tag         uint8
    name_index  uint16
}

// field ref
// method ref
// interface method ref
type cp_ref_info struct {
    tag                 uint8
    class_index         uint16
    name_and_type_index uint16
}

type cp_string_info struct {
    tag             uint8
    string_index    uint16
}

// integer
// float
type cp_u4_info struct {
    tag     uint8
    bytes   uint32
}

// long
// double
type cp_u8_info struct {
    tag         uint8
    high_bytes  uint32
    low_bytes   uint32
}

type field_info struct {
    access_flags        uint16
    name_index          uint16
    descriptor_index    uint16
    attributes_count    uint16
    attributes          []attribute_info
}

type method_info struct {

}

type attribute_info struct {
    attribute_name_index    uint16
    attribute_length        uint32
    info                    []uint8
}

type vmclass struct {
    magic               uint32
    minor_version       uint16
    major_version       uint16
    constant_pool_count uint16
    constant_pool       []cp_info
    access_flags        uint16
    this_class          uint16
    super_class         uint16
    interfaces_count    uint16
    interfaces          []uint16
    fields_count        uint16
    fields              []field_info
    methods_count       uint16
    methods             []method_info
    attributes_count    uint16
    attributes          []attribute_info
}

type decoder struct {
    file    io.Reader
    bo      binary.ByteOrder
    mc      *vmclass
}

func (d *decoder) readMagic() {
    binary.Read(d.file, d.bo, &(d.mc.magic))
}

func (d *decoder) readVersion() {
    binary.Read(d.file, d.bo, &(d.mc.minor_version))
    binary.Read(d.file, d.bo, &(d.mc.major_version))
}

func (d *decoder) readConstantPool() {
    binary.Read(d.file, d.bo, &(d.mc.constant_pool_count))
    d.mc.constant_pool = make([]cp_info, d.mc.constant_pool_count)
    fmt.Printf("cp count=%d\n", d.mc.constant_pool_count-1) // -1 just for skipping 0
    for i := uint16(1); i < d.mc.constant_pool_count; i++ {
        var tag uint8
        binary.Read(d.file, d.bo, &tag)
        switch tag {
            case CONST_Class:
                info := make([]byte, 2)
                binary.Read(d.file, d.bo, info)
                d.mc.constant_pool[i] = cp_info{tag: tag, info: info}

            case CONST_FieldRef:    fallthrough
            case CONST_MethodRef:   fallthrough
            case CONST_InterfaceMethodRef:
                info := make([]byte, 4)
                binary.Read(d.file, d.bo, info)
                d.mc.constant_pool[i] = cp_info{tag: tag, info: info}

            case CONST_String:
                info := make([]byte, 2)
                binary.Read(d.file, d.bo, info)
                d.mc.constant_pool[i] = cp_info{tag: tag, info: info}

            case CONST_Integer:     fallthrough
            case CONST_Float:
                info := make([]byte, 4)
                binary.Read(d.file, d.bo, info)
                d.mc.constant_pool[i] = cp_info{tag: tag, info: info}

            case CONST_Long:        fallthrough
            case CONST_Double:
                info := make([]byte, 8)
                binary.Read(d.file, d.bo, info)
                d.mc.constant_pool[i] = cp_info{tag: tag, info: info}

            case CONST_NameAndType:
                info := make([]byte, 4)
                binary.Read(d.file, d.bo, info)
                d.mc.constant_pool[i] = cp_info{tag: tag, info: info}

            case CONST_Utf8:
                length := make([]byte, 2)
                binary.Read(d.file, d.bo, length)
                len := d.bo.Uint16(length)
                info := make([]byte, len + 2)
                copy(info, length)
                binary.Read(d.file, d.bo, info[2:])
                d.mc.constant_pool[i] = cp_info{tag: tag, info: info}
        }
        // fmt.Printf("cp[%d]: %v\n", i, mc.constant_pool[i])
    }    
}

func (d *decoder) readInterface() {
    f  := d.file
    bo := d.bo
    mc := d.mc
    binary.Read(f, bo, &mc.access_flags)
    binary.Read(f, bo, &mc.this_class)
    binary.Read(f, bo, &mc.super_class)
    binary.Read(f, bo, &mc.interfaces_count)
    if(d.mc.interfaces_count > 0) {
        d.mc.interfaces = make([]uint16, mc.interfaces_count)
        binary.Read(f, bo, &mc.interfaces)
    }    
}

func (d *decoder) readFields() {
    f  :=  d.file
    bo := d.bo
    mc := d.mc
    binary.Read(f, bo, &mc.fields_count)
    if(mc.fields_count > 0) {
        mc.fields = make([]field_info, mc.fields_count)
        for i := uint16(0); i < mc.fields_count; i++ {
            var fi field_info
            binary.Read(f, bo, &fi.access_flags)
            binary.Read(f, bo, &fi.name_index)
            binary.Read(f, bo, &fi.descriptor_index)
            binary.Read(f, bo, &fi.attributes_count)
            if(fi.attributes_count > 0) {
                fi.attributes = make([]attribute_info, fi.attributes_count)
                for j := uint16(0); j < fi.attributes_count; j++ {
                    var name_index uint16
                    var length uint32
                    binary.Read(f, bo, &name_index)
                    binary.Read(f, bo, &length)                    
                    info := make([]uint8, length)
                    binary.Read(f, bo, &info)
                    fi.attributes[j] = attribute_info {
                        attribute_name_index: name_index,
                        attribute_length: length,
                        info: info,
                    }
                }
            }
            mc.fields[i] = fi
        }
    }
    
}

func readClass(fileName string) (mc vmclass) {
    f,_ := os.Open(fileName, os.O_RDONLY, 0666)
    defer f.Close()

    d := decoder { file: f, bo: binary.BigEndian, mc: &mc }
    d.readMagic()
    d.readVersion()
    d.readConstantPool()
    d.readInterface()
    d.readFields()

    return
}

func (mc vmclass) eachField(f func(fi field_info, name, desc string)) {
    cp := mc.constant_pool
    for i := uint16(0); i < mc.fields_count; i++ {
       fi := mc.fields[i]
       cp1 := cp[fi.name_index]
       cp2 := cp[fi.descriptor_index]
       f(fi, string(cp1.info[2:]), string(cp2.info[2:]))
    }
}

func main() {
    bo := binary.BigEndian
    className := os.Args[1]
    classFile := className + ".class"
    fmt.Printf("%s\n", className)
    mc := readClass(classFile)
    fmt.Printf("%x\n", mc.magic)
    fmt.Printf("%d.%d\n", mc.major_version, mc.minor_version)
    cp := mc.constant_pool
    for i := uint16(1); i < mc.constant_pool_count; i++ {
        cp := mc.constant_pool[i]
        if(cp.tag == CONST_Utf8) {
            fmt.Printf("%s\n", string(cp.info[2:]))
        } else {
            fmt.Printf("%v\n", cp)
        }
    }
    fmt.Printf("class : %s\n", string(cp[bo.Uint16(cp[mc.this_class ].info)].info[2:]))
    fmt.Printf("super : %s\n", string(cp[bo.Uint16(cp[mc.super_class].info)].info[2:]))
    mc.eachField(func(fi field_info, name, desc string) {
        fmt.Printf("fi: %v %s:%s\n", fi, name, desc)
    })
}