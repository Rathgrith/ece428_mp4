import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

import java.io.IOException;

public class InterconneMapper extends Mapper<Object, Text, Text, IntWritable> {
    private final static IntWritable one = new IntWritable(1);
    private Text detectionType = new Text();

    public void map(Object key, Text value, Context context) throws IOException, InterruptedException {
        String[] parts = value.toString().split(",");
        String interconneType = context.getConfiguration().get("interconneType");
        if (parts[interconneIndex].equals(interconneType)) {
            detectionType.set(parts[detectionIndex]);
            context.write(detectionType, one);
        }
    }
}
