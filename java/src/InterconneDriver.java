import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class InterconneDriver {
    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();
        conf.set("interconneType", args[0]);
        conf.setInt("totalCount", Integer.parseInt(args[1]));

        Job job = Job.getInstance(conf, "Interconne Count");
        job.setJarByClass(InterconneDriver.class);
        job.setMapperClass(InterconneMapper.class);
        job.setCombinerClass(InterconneReducer.class);
        job.setReducerClass(InterconneReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(IntWritable.class);

        FileInputFormat.addInputPath(job, new Path(args[2]));
        FileOutputFormat.setOutputPath(job, new Path(args[3]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
